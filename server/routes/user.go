package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	db *mongo.Database
}

func InitUserRoutes(db *mongo.Database, mux *http.ServeMux) {
	// Initialize the UserController with the database client
	userController := &UserController{db: db}

	mux.HandleFunc("/register", userController.Register)
	mux.HandleFunc("/login", userController.Login)
	// Use withAuth middleware for protected routes
	mux.Handle("/profile", withAuth(http.HandlerFunc(userController.GetProfile)))
}

func (ctlr UserController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Save the user to the database
	collection := ctlr.db.Collection("users")
	_, err = collection.InsertOne(context.Background(), map[string]string{
		"email":    req.Email,
		"password": string(hashedPassword),
	})
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User registered successfully")
}

func (ctlr UserController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the user exists
	collection := ctlr.db.Collection("users")
	var user struct {
		Email    string `bson:"email"`
		Password string `bson:"password"`
	}
	err := collection.FindOne(context.Background(), map[string]string{"email": req.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"email": req.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (ctlr UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Email not found in context", http.StatusUnauthorized)
		return
	}

	// Retrieve the user's profile from the database
	collection := ctlr.db.Collection("users")
	var user struct {
		Email string `bson:"email"`
	}
	err := collection.FindOne(context.Background(), map[string]string{"email": email}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
