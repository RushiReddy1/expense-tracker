package main

import (
	"fmt"
	"net/http"

	"expense-tracker-backend/config"
	"expense-tracker-backend/handlers"
)

func main() {
	config.ConnectDB()

	http.HandleFunc("/expense", handlers.CreateExpense)
	http.HandleFunc("/expenses", handlers.GetExpenses)
	http.HandleFunc("/expense/delete", handlers.DeleteExpense)
	http.HandleFunc("/expense/update", handlers.UpdateExpense)


	http.HandleFunc("/transaction", handlers.CreateTransaction)
	http.HandleFunc("/transactions", handlers.GetTransactions)
	http.HandleFunc("/summary", handlers.GetSummary)

	fmt.Println("🚀 Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}