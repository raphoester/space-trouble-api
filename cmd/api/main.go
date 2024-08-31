package main

import (
	"net"

	bookingsv1 "github.com/raphoester/space-trouble-api/generated/proto/bookings/v1"
	"github.com/raphoester/space-trouble-api/internal/domain/commands/book_ticket"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/primary/controller"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/hardcoded_destination_registry"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/hardcoded_launchpad_registry"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/inmemory_bookings_storage"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/inmemory_competitor_flights_provider"
	"github.com/raphoester/space-trouble-api/internal/queries/get_all_bookings/inmemory_bookings_getter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	bookingsRepo := inmemory_bookings_storage.New()
	competitorFlightsProvider := inmemory_competitor_flights_provider.New()
	launchpadRegistry := hardcoded_launchpad_registry.New()
	destinationRegistry := hardcoded_destination_registry.New()
	ticketBooker := book_ticket.NewTicketBooker(bookingsRepo, competitorFlightsProvider,
		launchpadRegistry, destinationRegistry)

	bookingsGetter := inmemory_bookings_getter.New(bookingsRepo)
	ctr := controller.New(ticketBooker, bookingsGetter)

	server := grpc.NewServer()
	bookingsv1.RegisterBookingsServiceServer(server, ctr)
	reflection.Register(server)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}
