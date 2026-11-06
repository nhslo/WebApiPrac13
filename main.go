package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

var database *sql.DB

func send(responseWriter http.ResponseWriter, status int, value any) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)
	json.NewEncoder(responseWriter).Encode(value)
}

func idFrom(responseWriter http.ResponseWriter, request *http.Request) (int, bool) {
	id, operationError := strconv.Atoi(request.PathValue("id"))
	if operationError != nil || id <= 0 {
		http.Error(responseWriter, "Invalid ID", 400)
		return 0, false
	}
	return id, true
}

func serverError(responseWriter http.ResponseWriter, operationError error) {
	log.Println(operationError)
	http.Error(responseWriter, "Database error", 500)
}

// Получаем товары из PostgreSQL с необязательным фильтром по категории.
func list(responseWriter http.ResponseWriter, request *http.Request) {
	category := request.URL.Query().Get("category")
	rows, operationError := database.Query("SELECT id,name,category,price,quantity FROM products WHERE ($1 = '' OR category = $1) ORDER BY id", category)
	if operationError != nil {
		serverError(responseWriter, operationError)
		return
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var product Product
		if operationError = rows.Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.Quantity); operationError != nil {
			serverError(responseWriter, operationError)
			return
		}
		items = append(items, product)
	}
	if operationError = rows.Err(); operationError != nil {
		serverError(responseWriter, operationError)
		return
	}
	send(responseWriter, 200, items)
}

func get(responseWriter http.ResponseWriter, request *http.Request) {
	id, success := idFrom(responseWriter, request)
	if !success {
		return
	}
	var product Product
	operationError := database.QueryRow("SELECT id,name,category,price,quantity FROM products WHERE id=$1", id).Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.Quantity)
	if errors.Is(operationError, sql.ErrNoRows) {
		http.NotFound(responseWriter, request)
		return
	}
	if operationError != nil {
		serverError(responseWriter, operationError)
		return
	}
	send(responseWriter, 200, product)
}

// Проверяем данные и добавляем товар параметризованным SQL-запросом.
func create(responseWriter http.ResponseWriter, request *http.Request) {
	var product Product
	if json.NewDecoder(request.Body).Decode(&product) != nil || strings.TrimSpace(product.Name) == "" || strings.TrimSpace(product.Category) == "" || product.Price < 0 || product.Quantity < 0 {
		http.Error(responseWriter, "Invalid product", 400)
		return
	}
	operationError := database.QueryRow("INSERT INTO products(name,category,price,quantity) VALUES($1,$2,$3,$4) RETURNING id", product.Name, product.Category, product.Price, product.Quantity).Scan(&product.ID)
	if operationError != nil {
		serverError(responseWriter, operationError)
		return
	}
	send(responseWriter, 201, product)
}

func update(responseWriter http.ResponseWriter, request *http.Request) {
	id, success := idFrom(responseWriter, request)
	if !success {
		return
	}
	var product struct {
		Price    *float64 `json:"price"`
		Quantity *int     `json:"quantity"`
	}
	if json.NewDecoder(request.Body).Decode(&product) != nil || product.Price == nil || product.Quantity == nil || *product.Price < 0 || *product.Quantity < 0 {
		http.Error(responseWriter, "Invalid product", 400)
		return
	}
	result, operationError := database.Exec("UPDATE products SET price=$1,quantity=$2 WHERE id=$3", *product.Price, *product.Quantity, id)
	if operationError != nil {
		serverError(responseWriter, operationError)
		return
	}
	numberValue, _ := result.RowsAffected()
	if numberValue == 0 {
		http.NotFound(responseWriter, request)
		return
	}
	get(responseWriter, request)
}

func remove(responseWriter http.ResponseWriter, request *http.Request) {
	id, success := idFrom(responseWriter, request)
	if !success {
		return
	}
	result, operationError := database.Exec("DELETE FROM products WHERE id=$1", id)
	if operationError != nil {
		serverError(responseWriter, operationError)
		return
	}
	numberValue, _ := result.RowsAffected()
	if numberValue == 0 {
		http.NotFound(responseWriter, request)
		return
	}
	responseWriter.WriteHeader(204)
}

func main() {
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		log.Fatal("Set DATABASE_URL")
	}
	var operationError error
	database, operationError = sql.Open("pgx", connectionString)
	if operationError != nil {
		log.Fatal(operationError)
	}
	defer database.Close()
	if operationError = database.Ping(); operationError != nil {
		log.Fatal(operationError)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", list)
	mux.HandleFunc("GET /products/{id}", get)
	mux.HandleFunc("POST /products", create)
	mux.HandleFunc("PUT /products/{id}", update)
	mux.HandleFunc("DELETE /products/{id}", remove)
	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
