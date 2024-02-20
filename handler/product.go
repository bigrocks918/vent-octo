package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"octo-api/model"
	"octo-api/store"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetProducts godoc
// @Summary Get all products
// @Description Retrieves all products, with the option to filter by pricing mode
// @Tags product
// @Accept  json
// @Produce  json
// @Param   Capability header string false "Capability to filter by pricing mode"
// @Success 200 {array} model.ProductPayload_Rs_NonPricing "Success - Return all products in non-pricing mode"
// @Success 200 {array} model.Product "Success - Return all products in pricing mode"
// @Failure 500 {string} string "Internal Server Error"
// @Router /products [get]
func GetProducts(w http.ResponseWriter, r *http.Request) {
	capHeader := r.Header.Get("Capability")
	// Check if pricing mode
	isExt := (strings.ToLower(capHeader) == "pricing")

	database := store.ConnectToDB()
	defer database.Close()

	// Get the Whole Product Data from DB
	products, err := store.GetProductsFromDB(database)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Prepare Output data according to mode
	if isExt { // pricing mode
		var productsOutputs []model.ProductPayload_Rs_Pricing
		for _, product := range products {
			productsOutputs = append(productsOutputs, model.ProductPayload_Rs_Pricing{
				Id:       product.ID,
				Name:     product.Name,
				Capacity: product.Capacity,
				Price:    product.Price,
				Currency: product.Currency,
			})
		}
		json.NewEncoder(w).Encode(productsOutputs)
	} else { // Non-Pricing mode
		var productsOutputs []model.ProductPayload_Rs_NonPricing
		for _, product := range products {
			productsOutputs = append(productsOutputs, model.ProductPayload_Rs_NonPricing{
				Id:       product.ID,
				Name:     product.Name,
				Capacity: product.Capacity,
			})
		}
		json.NewEncoder(w).Encode(productsOutputs)
	}
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Fetches a product by its ID, with the option to filter by pricing mode
// @Tags product
// @Accept  json
// @Produce  json
// @Param   id path string true "Product ID"
// @Param   Capability header string false "Capability to filter by pricing mode"
// @Success 200 {object} model.ProductPayload_Rs_NonPricing "Success - Return product in non-pricing mode"
// @Success 200 {object} model.Product "Success - Return product in pricing mode"
// @Failure 404 {string} string "Product not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /products/{id} [get]
func GetProduct(w http.ResponseWriter, r *http.Request) {

	// Get ProductID
	vars := mux.Vars(r)
	productId := vars["id"]

	capHeader := r.Header.Get("Capability")
	// Check if pricing mode
	isExt := (strings.ToLower(capHeader) == "pricing")

	database := store.ConnectToDB()
	defer database.Close()

	// Get Product with certain ID
	product, err := store.GetProductFromDB(database, productId)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Prepare Output data according to mode
	if isExt { // Pricing mode
		outputProduct := model.ProductPayload_Rs_Pricing{
			Id:       product.ID,
			Name:     product.Name,
			Capacity: product.Capacity,
			Price:    product.Price,
			Currency: product.Currency,
		}
		json.NewEncoder(w).Encode(outputProduct)
	} else { // Non-Pricing mode
		outputProduct := model.ProductPayload_Rs_NonPricing{
			Id:       product.ID,
			Name:     product.Name,
			Capacity: product.Capacity,
		}
		json.NewEncoder(w).Encode(outputProduct)
	}
}

// AddProduct godoc
// @Summary Add a product
// @Description Adds a new product to the database
// @Tags product
// @Accept  json
// @Produce  json
// @Param   ProductPayload_Rq body model.ProductPayload_Rq true "Request Payload for Adding a Product"
// @Success 201 {string} string "successfully created"
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal Server Error"
// @Router /products/add [post]
func AddProduct(w http.ResponseWriter, r *http.Request) {
	// Decode Product Data from request
	var product_schema model.ProductPayload_Rq
	if err := json.NewDecoder(r.Body).Decode(&product_schema); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	database := store.ConnectToDB()
	defer database.Close()

	if len(product_schema.Currency) == 0 {
		// Set the default currency type as USD if it is not mentioned in payload
		product_schema.Currency = "USD"
	}

	// Add Product to DB
	err := store.InsertProductIntoDB(database, model.Product{
		ID:       uuid.NewString(),
		Name:     product_schema.Name,
		Capacity: product_schema.Capacity,
		Price:    product_schema.Price,
		Currency: product_schema.Currency,
	})
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("successfully created")
}
