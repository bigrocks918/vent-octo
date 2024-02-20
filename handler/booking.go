package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"octo-api/helper"
	"octo-api/model"
	"octo-api/store"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// PostBooking godoc
// @Summary Post a booking
// @Description Creates a new booking and updates the availability accordingly
// @Tags booking
// @Accept  json
// @Produce  json
// @Param   BookingPayload_Rq body model.BookingPayload_Rq true "Request Payload for Posting a Booking"
// @Success 201 {object} model.Booking "Booking successfully created"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal Server Error"
// @Router /bookings/post [post]
func PostBooking(w http.ResponseWriter, r *http.Request) {

	// Decode Booking info from request
	var bookingSchema model.BookingPayload_Rq
	if err := json.NewDecoder(r.Body).Decode(&bookingSchema); err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	database := store.ConnectToDB()
	defer database.Close()

	// Check if availabilityId is Valid & Check Price and Currency
	// Get Availability with certain AvailabilityID
	availability, err := store.GetAvailabilityByIdFromDB(database, bookingSchema.AvailabilityId)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
		http.Error(w, "Invalid AvailabilityID", http.StatusBadRequest)
		return
	}
	// Get Product information with certain ProductID
	product, err := store.GetProductFromDB(database, availability.ProductId)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
		http.Error(w, "Internal DB Error", http.StatusInternalServerError)
		return
	}

	var booking model.Booking

	booking.AvailabilityId = bookingSchema.AvailabilityId
	booking.Units = bookingSchema.Units

	// Generate a unique ID for the new booking
	booking.ID = uuid.New().String()
	booking.Status = "RESERVED"

	// Calculate Price
	// Currency check
	// LOGIC : If currency of product and availability is not matched together, set booking currency as USD. Else, continue with the same currency
	if product.Currency != availability.Currency {
		// Set USD as currency unit for booking
		booking.Currency = "USD"
		booking.Price = 0
		if strings.ToLower(product.Currency) == "usd" {
			booking.Price = booking.Price + product.Price*float64(bookingSchema.Units)
		} else {
			// Make conversion
			converted_amount, err := helper.Rate_Convert(product.Currency, booking.Currency, product.Price)
			if err != nil {
				// log.Fatal(err)
				fmt.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			booking.Price = booking.Price + converted_amount*float64(bookingSchema.Units)
		}

		if strings.ToLower(availability.Currency) == "usd" {
			booking.Price = booking.Price + availability.Price*float64(bookingSchema.Units)
		} else {
			// Make conversion
			converted_amount, err := helper.Rate_Convert(availability.Currency, booking.Currency, availability.Price)
			if err != nil {
				// log.Fatal(err)
				fmt.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			booking.Price = booking.Price + converted_amount*float64(bookingSchema.Units)
		}
	} else {
		// Set the booking currency type as product's and availability's currency type
		booking.Currency = product.Currency
		booking.Price = (product.Price + availability.Price) * float64(bookingSchema.Units)
	}

	if err := store.CreateBooking(database, booking); err != nil {
		// log.Fatal(err)
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

// GetAllBookings godoc
// @Summary Get all bookings
// @Description Retrieves all bookings, with the option to filter by pricing mode
// @Tags booking
// @Accept  json
// @Produce  json
// @Param   Capability header string false "Capability to filter by pricing mode"
// @Success 200 {array} model.BookingPayload_Rs_NonPricing "Success - Return all bookings in non-pricing mode"
// @Success 200 {array} model.Booking "Success - Return all bookings in pricing mode"
// @Failure 500 {string} string "Internal Server Error"
// @Router /bookings/all [get]
func GetAllBookings(w http.ResponseWriter, r *http.Request) {

	capHeader := r.Header.Get("Capability")
	// Check if pricing mode
	isExt := (strings.ToLower(capHeader) == "pricing")

	database := store.ConnectToDB()
	defer database.Close()

	// Get All Booking lists
	bookings, err := store.GetAllBookings(database)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Prepare Out data according to the mode
	if !isExt {
		// non-pricing mode : remove all data related to pricing in output
		var nonPricingBookings []model.BookingPayload_Rs_NonPricing
		for _, booking := range bookings {
			var nonPricingBookingUnits []model.BookingUnitPayload_Rs_NonPricing
			for _, booking_unit := range booking.Units {
				nonPricingBookingUnits = append(nonPricingBookingUnits, model.BookingUnitPayload_Rs_NonPricing{
					ID:        booking_unit.ID,
					BookingId: booking_unit.BookingId,
				})
			}

			nonPricingBookings = append(nonPricingBookings, model.BookingPayload_Rs_NonPricing{
				ID:             booking.ID,
				Status:         booking.Status,
				AvailabilityId: booking.AvailabilityId,
				Units:          nonPricingBookingUnits,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(nonPricingBookings)
		return
	}

	// Pricing Mode
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookings)
}

// GetBooking godoc
// @Summary Get a booking by ID
// @Description Fetches a booking by its ID, with the option to filter by pricing mode
// @Tags booking
// @Accept  json
// @Produce  json
// @Param   id path string true "Booking ID"
// @Param   Capability header string false "Capability to filter by pricing mode"
// @Success 200 {object} model.BookingPayload_Rs_NonPricing "Success - Return booking in non-pricing mode"
// @Success 200 {object} model.Booking "Success - Return booking in pricing mode"
// @Failure 404 {string} string "Booking not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /bookings/{id} [get]
func GetBooking(w http.ResponseWriter, r *http.Request) {

	capHeader := r.Header.Get("Capability")
	// Check if pricing mode
	isExt := (strings.ToLower(capHeader) == "pricing")

	vars := mux.Vars(r)
	bookingID := vars["id"]

	database := store.ConnectToDB()
	defer database.Close()

	// Get booking info with Id
	booking, err := store.GetBookingByID(database, bookingID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Prepare Out data according to the mode
	if !isExt {
		// non-pricing mode : remove all data related to pricing in output

		var nonPricingBookingUnits []model.BookingUnitPayload_Rs_NonPricing
		for _, booking_unit := range booking.Units {
			nonPricingBookingUnits = append(nonPricingBookingUnits, model.BookingUnitPayload_Rs_NonPricing{
				ID:        booking_unit.ID,
				BookingId: booking_unit.BookingId,
			})
		}

		nonPricingBooking := model.BookingPayload_Rs_NonPricing{
			ID:             booking.ID,
			Status:         booking.Status,
			AvailabilityId: booking.AvailabilityId,
			Units:          nonPricingBookingUnits,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(nonPricingBooking)
		return
	}

	// Pricing Mode
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

// ConfirmBooking godoc
// @Summary Confirm a booking
// @Description Confirms a booking by its ID
// @Tags booking
// @Accept  json
// @Produce  json
// @Param   id path string true "Booking ID to confirm"
// @Success 200 {string} string "Booking confirmed successfully"
// @Failure 404 {string} string "Booking not found after confirmation"
// @Failure 500 {string} string "Internal Server Error"
// @Router /bookings/confirm/{id} [put]
func ConfirmBooking(w http.ResponseWriter, r *http.Request) {

	// Get ID from Request URL
	vars := mux.Vars(r)
	bookingID := vars["id"]

	database := store.ConnectToDB()
	defer database.Close()

	// Confirm Booking with id
	if err := store.ConfirmBooking(database, bookingID); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get booking with ID
	booking, err := store.GetBookingByID(database, bookingID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Booking not found after confirmation", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}
