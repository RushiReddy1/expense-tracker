package handlers

import (
	"encoding/json"
	"expense-tracker-backend/config"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)


// -------------------- CREATE TRANSACTION --------------------
type Transaction struct {
	Title     string    `json:"title" bson:"title"`
	Amount    float64   `json:"amount" bson:"amount"`
	Type      string    `json:"type" bson:"type"`
	Source    string    `json:"source" bson:"source"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var tx Transaction

	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// Add timestamp (stored correctly as BSON Date)
	tx.CreatedAt = time.Now()

	collection := config.Client.
		Database("expense_tracker").
		Collection("transactions")

	_, err = collection.InsertOne(r.Context(), tx)
	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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

	var transactions []Transaction

	for cursor.Next(r.Context()) {
		var tx Transaction

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
			"created_at": bson.M{
				"$gte": startTime,
				"$lt":  endTime,
			},
		}
	} else {
		filter = bson.M{}
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
		var tx Transaction
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