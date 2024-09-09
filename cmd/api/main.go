package main

import (
	"fmt"
	"net"
	"os"

	bookingsv1 "github.com/raphoester/space-trouble-api/generated/proto/bookings/v1"
	"github.com/raphoester/space-trouble-api/internal/domain/commands/book_ticket"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/primary/controller"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/hardcoded_destination_registry"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/hardcoded_launchpad_registry"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/psql_bookings_storage"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/spacex_competitor_flights_provider"
	"github.com/raphoester/space-trouble-api/internal/pkg/postgres"
	"github.com/raphoester/space-trouble-api/internal/queries/get_all_bookings/psql_bookings_getter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO: manage app DI sequence with a separate package

func main() {
	// TODO: replace with viper (for the demo this is fine)
	pgDSN := os.Getenv("POSTGRES_DSN")
	migrationsPath := os.Getenv("MIGRATIONS_PATH")

	pg, err := postgres.New(pgDSN)
	if err != nil {
		panic(err)
	}

	// TODO: add condition to run migrations only in dev environment
	if err := pg.Migrate(migrationsPath); err != nil {
		panic(err)
	}

	bookingsRepo := psql_bookings_storage.New(pg)
	competitorFlightsProvider := spacex_competitor_flights_provider.New()
	launchpadRegistry := hardcoded_launchpad_registry.New()
	destinationRegistry := hardcoded_destination_registry.New()
	ticketBooker := book_ticket.NewTicketBooker(bookingsRepo, competitorFlightsProvider,
		launchpadRegistry, destinationRegistry)

	bookingsGetter := psql_bookings_getter.New(pg)
	ctr := controller.New(ticketBooker, bookingsGetter)

	server := grpc.NewServer()
	bookingsv1.RegisterBookingsServiceServer(server, ctr)
	reflection.Register(server)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	// TODO: add proper logging
	fmt.Printf("Listening on %s\n", listener.Addr().String())
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
