package model

import "time"

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Capacity int     `json:"capacity"`
	Price    float64 `json:"price,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

type Availability struct {
	ID        string    `json:"id"`
	LocalDate time.Time `json:"localDate"`
	Status    string    `json:"status"`
	ProductId string    `json:"productId"`
	Vacancies int       `json:"vacancies"`
	Available bool      `json:"available"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
}

type AvailabilityShow struct {
	ID          string    `json:"id"`
	LocalDate   time.Time `json:"localDate"`
	Status      string    `json:"status"`
	ProductName string    `json:"productName"`
	Vacancies   int       `json:"vacancies"`
	Available   bool      `json:"available"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
}

type Booking struct {
	ID             string  `json:"id"`
	Status         string  `json:"status"`
	AvailabilityId string  `json:"availabilityId"`
	Units          int     `json:"units"`
	Price          float64 `json:"price"`
	Currency       string  `json:"currency"`
}

type BookingUnit struct {
	ID        string  `json:"id"`
	BookingId string  `json:"bookingId"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
}
