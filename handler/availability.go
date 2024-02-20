package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"octo-api/model"
	"octo-api/store"
	"strings"
	"time"
)

// GetAvailabilities godoc
// @Summary Get availabilities
// @Description Get availabilities by single date or date range
// @Tags availability
// @Accept  json
// @Produce  json
// @Param   Capability header string false "Capability"
// @Param   AvailabilityPayload_Rq body model.AvailabilityPayload_Rq true "Request Payload"
// @Success 200 {object} []model.AvailabilityPayload_Rs_Pricing "Success"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal Server Error"
// @Router /availabilities [get]
func GetAvailabilities(w http.ResponseWriter, r *http.Request) {

	capHeader := r.Header.Get("Capability")
	// Check if pricing mode
	isExt := (strings.ToLower(capHeader) == "pricing")

	// Decode date information from request
	var req model.AvailabilityPayload_Rq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var startDate, endDate time.Time
	var err error

	if req.LocalDate != "" {
		// Single date query
		startDate, err = time.Parse("2006-01-02", req.LocalDate)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Invalid localDate format. Please use YYYY-MM-DD.", http.StatusBadRequest)
			return
		}
		endDate = startDate // Single date, so start and end are the same
	} else {
		// Date range query
		startDate, err = time.Parse("2006-01-02", req.LocalDateStart)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Invalid localDateStart format. Please use YYYY-MM-DD.", http.StatusBadRequest)
			return
		}
		endDate, err = time.Parse("2006-01-02", req.LocalDateEnd)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Invalid localDateEnd format. Please use YYYY-MM-DD.", http.StatusBadRequest)
			return
		}
	}

	database := store.ConnectToDB()
	defer database.Close()

	// Get Availability Data
	availabilities, err := store.GetAvailabilitiesFromDB(database, startDate, endDate)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Prepare output data according to mode
	if isExt { // Pricing mode
		var availabilityOutputs []model.AvailabilityPayload_Rs_Pricing
		for _, availability := range availabilities {
			availabilityOutputs = append(
				availabilityOutputs,
				model.AvailabilityPayload_Rs_Pricing{
					Id:          availability.ID,
					LocalDate:   availability.LocalDate,
					Status:      availability.Status,
					ProductName: availability.ProductName,
					Vacancies:   availability.Vacancies,
					Available:   availability.Available,
					Price:       availability.Price,
					Currency:    availability.Currency,
				},
			)
		}
		json.NewEncoder(w).Encode(availabilityOutputs)
	} else { // Non-Pricing mode
		var availabilityOutputs []model.AvailabilityPayload_Rs_NonPricing
		for _, availability := range availabilities {
			availabilityOutputs = append(
				availabilityOutputs,
				model.AvailabilityPayload_Rs_NonPricing{
					Id:          availability.ID,
					LocalDate:   availability.LocalDate,
					Status:      availability.Status,
					ProductName: availability.ProductName,
					Vacancies:   availability.Vacancies,
					Available:   availability.Available,
				},
			)
		}
		json.NewEncoder(w).Encode(availabilityOutputs)
	}
}

// AddAvailabilities godoc
// @Summary Add availabilities
// @Description Adds availabilities for a product within a single date or date range
// @Tags availability
// @Accept  json
// @Produce  json
// @Param   AvailabilityNewPayload_Rq body model.AvailabilityNewPayload_Rq true "Request Payload for Adding Availabilities"
// @Success 201 {string} string "successfully added"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal Server Error"
// @Router /availabilities/add [post]
func AddAvailabilities(w http.ResponseWriter, r *http.Request) {

	var req model.AvailabilityNewPayload_Rq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	database := store.ConnectToDB()
	defer database.Close()

	var startDate, endDate time.Time
	var err error

	if req.LocalDate != "" {
		// Single date query
		startDate, err = time.Parse("2006-01-02", req.LocalDate)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Invalid localDate format. Please use YYYY-MM-DD.", http.StatusBadRequest)
			return
		}
		endDate = startDate // Single date, so start and end are the same
	} else {
		// Date range query
		startDate, err = time.Parse("2006-01-02", req.LocalDateStart)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Invalid localDateStart format. Please use YYYY-MM-DD.", http.StatusBadRequest)
			return
		}
		endDate, err = time.Parse("2006-01-02", req.LocalDateEnd)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Invalid localDateEnd format. Please use YYYY-MM-DD.", http.StatusBadRequest)
			return
		}
	}

	if len(req.Currency) == 0 {
		// Set default currency type as USD if use didn't input currency information
		req.Currency = "USD"
	}

	err = store.AddAvailabilityIntoDB(database, req.ProductId, startDate, endDate, req.Price, req.Currency)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("successfully added")
}
