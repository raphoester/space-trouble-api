package bookings

import "github.com/raphoester/space-trouble-api/internal/pkg/birthday"

type ClientData struct {
	FirstName string
	LastName  string
	Gender    string
	Birthday  birthday.Date
}
