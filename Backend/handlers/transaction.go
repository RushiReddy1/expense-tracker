package handlers

import (
	"encoding/json"
	"expense-tracker-backend/config"
	"expense-tracker-backend/models"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// -------------------- CREATE TRANSACTION ---------


func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var tx models.Transaction
	if tx.Amount <= 0 {
	http.Error(w, "Invalid amount", http.StatusBadRequest)
	return
}

if tx.Type == "" {
	http.Error(w, "Type required", http.StatusBadRequest)
	return
}

json.NewDecoder(r.Body).Decode(&tx)


email := r.Context().Value("userEmail").(string)
tx.Email = email

// Add timestamp
tx.CreatedAt = time.Now()

	// Add timestamp (stored correctly as BSON Date)
	tx.CreatedAt = time.Now()

	collection := config.Client.
		Database("expense_tracker").
		Collection("transactions")

	_, err := collection.InsertOne(r.Context(), tx)
	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tx)
}

// -------------------- GET ALL TRANSACTIONS --------------------

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	collection := config.Client.
		Database("expense_tracker").
		Collection("transactions")

	cursor, err := collection.Find(r.Context(), bson.M{})
	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, "Failed to fetch", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var transactions []models.Transaction

	for cursor.Next(r.Context()) {
		var tx models.Transaction

		err := cursor.Decode(&tx)
		if err != nil {
			fmt.Println("Decode error:", err)
			continue
		}

		transactions = append(transactions, tx)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// -------------------- GET SUMMARY --------------------

func GetSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	email := r.Context().Value("userEmail").(string)
fmt.Println("Summary user:", email)

	// 1. Read query param
	monthParam := r.URL.Query().Get("month")

	var filter bson.M

	// 2. Build filter if month provided
	if monthParam != "" {
		startTime, err := time.Parse("2006-01", monthParam)
		if err != nil {
			http.Error(w, "Invalid month format", http.StatusBadRequest)
			return
		}

		endTime := startTime.AddDate(0, 1, 0)

		filter = bson.M{
			   "email": email, 
			"created_at": bson.M{
				"$gte": startTime,
				"$lt":  endTime,
			},
		}
	} else {
		filter = bson.M{
    "email": email,  
}
	}

	collection := config.Client.
		Database("expense_tracker").
		Collection("transactions")

	cursor, err := collection.Find(r.Context(), filter)
	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, "Failed to fetch", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	// 3. Totals
	var totalExpense float64
	var totalInvestment float64
	var totalTransfer float64

	// 4. Loop through filtered results
	for cursor.Next(r.Context()) {
		var tx models.Transaction
		cursor.Decode(&tx)

		amount := tx.Amount
		typeVal := tx.Type

		if typeVal == "expense" {
			totalExpense += amount
		} else if typeVal == "investment" {
			totalInvestment += amount
		} else if typeVal == "transfer" {
			totalTransfer += amount
		}
	}

	// 5. Optional toggle
	includeInvestments := r.URL.Query().Get("includeInvestments")

	var netOutflow float64
	if includeInvestments == "true" {
		netOutflow = totalExpense + totalInvestment
	} else {
		netOutflow = totalExpense
	}

	// 6. Response
	response := map[string]interface{}{
		"total_expenses":    totalExpense,
		"total_investments": totalInvestment,
		"total_transfers":   totalTransfer,
		"net_outflow":       netOutflow,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
