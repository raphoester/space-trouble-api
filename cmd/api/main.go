package main

import (
	"net"

	bookingsv1 "github.com/raphoester/space-trouble-api/generated/proto/bookings/v1"
	"github.com/raphoester/space-trouble-api/internal/domain/commands/book_ticket"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/primary/controller"
	"github.com/raphoester/space-trouble-api/internal/infrastructure/secondary/inmemory_bookings_storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	bookingsRepo := inmemory_bookings_storage.New()
	ticketBooker := book_ticket.NewTicketBooker(bookingsRepo)
	ctr := controller.New(ticketBooker)

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
