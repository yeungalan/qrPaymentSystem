package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type User struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Balance  float64 `json:"balance"`
}

type Transaction struct {
	ID        int       `json:"id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

var users = make(map[string]User)
var transactions []Transaction
var paymentIDCounter = 1

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/getBalance", getBalance).Methods("GET")
	r.HandleFunc("/getTransactionLog", getTransactionLog).Methods("GET")
	r.HandleFunc("/makeTransaction", makeTransaction).Methods("GET")
	r.HandleFunc("/register", register).Methods("POST")
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/getPaymentID", getPaymentID).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getBalance(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	user, ok := users[username]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	balance := user.Balance
	response := map[string]float64{"balance": balance}
	json.NewEncoder(w).Encode(response)
}

func getTransactionLog(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(transactions)
}

func makeTransaction(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	user, ok := users[username]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	paymentCode := r.URL.Query().Get("paymentCode")
	// Process paymentCode and update transaction log
	transactions = append(transactions, Transaction{
		ID:        len(transactions) + 1,
		Amount:    100.0, // Example amount
		Timestamp: time.Now(),
	})

	fmt.Fprint(w, "OK")
}

func register(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	users[user.Username] = user
	fmt.Fprint(w, "User registered")
}

func login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	storedUser, ok := users[user.Username]
	if !ok || storedUser.Password != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	fmt.Fprint(w, "Login successful")
}

func getPaymentID(w http.ResponseWriter, r *http.Request) {
	paymentID := fmt.Sprintf("PAY-%d", paymentIDCounter)
	paymentIDCounter++
	response := map[string]string{"paymentID": paymentID}
	json.NewEncoder(w).Encode(response)
}
