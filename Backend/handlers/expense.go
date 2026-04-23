package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"expense-tracker-backend/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
// Handler to create a new expense-POST /expense
func CreateExpense(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var expense map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	collection := config.Client.Database("expense_tracker").Collection("expenses")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, expense)
	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	fmt.Println("Saved:", expense)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Expense saved",
	})
}

// New handler to fetch all expenses-GET /expenses

func GetExpenses(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	collection := config.Client.Database("expense_tracker").Collection("expenses")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		http.Error(w, "Failed to fetch", http.StatusInternalServerError)
		return
	}

	var expenses []map[string]interface{}

	for cursor.Next(ctx) {
		var expense map[string]interface{}
		cursor.Decode(&expense)
		expenses = append(expenses, expense)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

//delete handler to delete an expense by ID-DELETE 

	func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	collection := config.Client.Database("expense_tracker").Collection("expenses")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, map[string]interface{}{"_id": objID})
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Deleted successfully",
	})
}

// New handler to delete an expense by ID-DELETE /expense?id=123 
func UpdateExpense(w http.ResponseWriter, r *http.Request) {

	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var updateData map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	collection := config.Client.Database("expense_tracker").Collection("expenses")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.UpdateOne(
		ctx,
		map[string]interface{}{"_id": objID},
		map[string]interface{}{"$set": updateData},
	)

	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Updated successfully",
	})
}