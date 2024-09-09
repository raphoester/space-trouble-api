package main

import (
	"flag"
	"fmt"
	"net"

	bookingsv1 "github.com/raphoester/space-trouble-api/generated/proto/bookings/v1"
	"github.com/raphoester/space-trouble-api/internal/domain/commands/book_ticket"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/primary/controller"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/hardcoded_destination_registry"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/hardcoded_launchpad_registry"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/psql_bookings_storage"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/spacex_competitor_flights_provider"
	"github.com/raphoester/space-trouble-api/internal/pkg/cfgutil"
	"github.com/raphoester/space-trouble-api/internal/pkg/postgres"
	"github.com/raphoester/space-trouble-api/internal/queries/get_all_bookings/psql_bookings_getter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO: manage app DI sequence with a separate package

type Config struct {
	PostgresDSN    string
	MigrationsPath string
	GRPCServer     struct {
		BindAddress      string
		EnableReflection bool
	}
}

func main() {
	c := flag.String("config", "./configs/sample.yaml", "path to config file")
	flag.Parse()

	cfg := Config{}
	if err := cfgutil.NewLoader(*c).Unmarshal(&cfg); err != nil {
		panic(err)
	}

	pg, err := postgres.New(cfg.PostgresDSN)
	if err != nil {
		panic(err)
	}

	// TODO: add condition to run migrations only in dev environment
	if err := pg.Migrate(cfg.MigrationsPath); err != nil {
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
	if cfg.GRPCServer.EnableReflection {
		reflection.Register(server)
	}

	listener, err := net.Listen("tcp", cfg.GRPCServer.BindAddress)
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
