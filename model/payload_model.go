package model

import "time"

type ProductPayload_Rq struct {
	Name     string  `json:"name"`
	Capacity int     `json:"capacity"`
	Price    float64 `json:"price,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

type ProductPayload_Rs_NonPricing struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type ProductPayload_Rs_Pricing struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Capacity int     `json:"capacity"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

type AvailabilityPayload_Rq struct {
	LocalDate      string `json:"localDate,omitempty"`
	LocalDateStart string `json:"localDateStart,omitempty"`
	LocalDateEnd   string `json:"localDateEnd,omitempty"`
}

type AvailabilityNewPayload_Rq struct {
	ProductId      string  `json:"productId"`
	LocalDate      string  `json:"localDate,omitempty"`
	LocalDateStart string  `json:"localDateStart,omitempty"`
	LocalDateEnd   string  `json:"localDateEnd,omitempty"`
	Price          float64 `json:"price,omitempty"`
	Currency       string  `json:"currency,omitempty"`
}

type AvailabilityPayload_Rs_NonPricing struct {
	Id          string    `json:"id"`
	LocalDate   time.Time `json:"localDate"`
	Status      string    `json:"status"`
	ProductName string    `json:"productName"`
	Vacancies   int       `json:"vacancies"`
	Available   bool      `json:"available"`
}

type AvailabilityPayload_Rs_Pricing struct {
	Id          string    `json:"id"`
	LocalDate   time.Time `json:"localDate"`
	Status      string    `json:"status"`
	ProductName string    `json:"productName"`
	Vacancies   int       `json:"vacancies"`
	Available   bool      `json:"available"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
}

type BookingPayload_Rq struct {
	ProductId      string `json:"productId,omitempty"`
	AvailabilityId string `json:"availabilityId"`
	Units          int    `json:"units"`
}

type BookingPayload_Rs struct {
	ID             string                  `json:"id"`
	Status         string                  `json:"status"`
	AvailabilityId string                  `json:"availabilityId"`
	Units          []BookingUnitPayload_Rs `json:"units"`
	Price          float64                 `json:"price"`
	Currency       string                  `json:"currency"`
}

type BookingUnitPayload_Rs struct {
	ID        string  `json:"id"`
	BookingId string  `json:"bookingId"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
}

type BookingPayload_Rs_NonPricing struct {
	ID             string                             `json:"id"`
	Status         string                             `json:"status"`
	AvailabilityId string                             `json:"availabilityId"`
	Units          []BookingUnitPayload_Rs_NonPricing `json:"units"`
}

type BookingUnitPayload_Rs_NonPricing struct {
	ID        string `json:"id"`
	BookingId string `json:"bookingId"`
}
