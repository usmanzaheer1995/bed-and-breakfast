package repository

import "github.com/usmanzaheer1995/bed-and-breakfast/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error
}
