package handlers

import (
	"encoding/json"
	"expense-tracker-backend/config"
	"expense-tracker-backend/models"
	"fmt"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
if err != nil {
	http.Error(w, "Error hashing password", http.StatusInternalServerError)
	return
}

user.Password = string(hashedPassword)
	collection := config.Client.
		Database("expense_tracker").
		Collection("users")

		var existingUser models.User

err = collection.FindOne(r.Context(), bson.M{"email": user.Email}).Decode(&existingUser)

if err == nil {
	// User found → email already exists
	http.Error(w, "Email already exists", http.StatusBadRequest)
	return
}
	_, err = collection.InsertOne(r.Context(), user)
	if err != nil {
		fmt.Println("DB ERROR:", err)
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
//login handler

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.User

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	collection := config.Client.
		Database("expense_tracker").
		Collection("users")

	var user models.User

	err = collection.FindOne(r.Context(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"email": user.Email,
	"exp":   time.Now().Add(time.Hour * 24).Unix(),
})

tokenString, err := token.SignedString([]byte("secret_key"))
if err != nil {
	http.Error(w, "Error generating token", http.StatusInternalServerError)
	return
}

	response := map[string]string{
	"token": tokenString,
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(response)
}
