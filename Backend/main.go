package main

import (
	"fmt"
	"net/http"

	"expense-tracker-backend/config"
	"expense-tracker-backend/handlers"
	"expense-tracker-backend/middleware"
)

func main() {
	config.ConnectDB()

	http.HandleFunc("/expense", handlers.CreateExpense)
	http.HandleFunc("/expenses", handlers.GetExpenses)
	http.HandleFunc("/expense/delete", handlers.DeleteExpense)
	http.HandleFunc("/expense/update", handlers.UpdateExpense)


	http.HandleFunc("/transaction", middleware.AuthMiddleware(handlers.CreateTransaction))
	http.HandleFunc("/transactions", middleware.AuthMiddleware(handlers.GetTransactions))
	http.HandleFunc("/summary", middleware.AuthMiddleware(handlers.GetSummary))
	http.HandleFunc("/signup", handlers.Signup)
	http.HandleFunc("/login", handlers.Login)

	fmt.Println("🚀 Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}