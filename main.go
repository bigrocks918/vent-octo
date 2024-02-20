package main

import (
	"net/http"
	"octo-api/handler"
	"octo-api/helper"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "octo-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {
	err := godotenv.Load()
	var apiKey string
	if err != nil {
		// default API
		apiKey = "cur_live_cAD8kz2VUuNH959jpUV3fukHiHOuo6RbHzHKnmPp"
	} else {
		apiKey = os.Getenv("CURRENCY_EXCHANGE_API_KEY")
	}
	helper.API_KEY = apiKey
}

func main() {
	r := mux.NewRouter()

	// Product routes
	r.HandleFunc("/products", handler.GetProducts).Methods("GET")
	r.HandleFunc("/products/new", handler.AddProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handler.GetProduct).Methods("GET")

	// Availability routes
	r.HandleFunc("/availability", handler.GetAvailabilities).Methods("GET")
	r.HandleFunc("/availability/add", handler.AddAvailabilities).Methods("POST")

	// Booking routes
	r.HandleFunc("/bookings", handler.PostBooking).Methods("POST")
	r.HandleFunc("/bookings/all", handler.GetAllBookings).Methods("GET")
	r.HandleFunc("/bookings/{id}", handler.GetBooking).Methods("GET")
	r.HandleFunc("/bookings/{id}/confirm", handler.ConfirmBooking).Methods("POST")

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	http.ListenAndServe(":8080", r)
}
