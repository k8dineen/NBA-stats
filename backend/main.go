package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"full_name"`
}

type Player struct {
	ID           int    `json:"id"`
	First        string `json:"first_name"`
	Last         string `json:"last_name"`
	Full         string `json:"full_name"`
	Position     string `json:"position"`
	Height       string `json:"height"`
	Weight       string `json:"weight"`
	Jersey       string `json:"jersey_number"`
	College      string `json:"college"`
	Draft_Year   int    `json:"draft_year"`
	Draft_Round  int    `json:"draft_round"`
	Draft_Number int    `json:"draft_number"`
	Team         Team   `json:"team"`
}

type PlayerResponse struct {
	Data []Player `json:"data"`
}

// Wrapper for player ID lookup response
type PlayerIdResponse struct {
	Data Player `json:"data"`
}

// Load API Key from .env
func loadAPIKey() string {
	godotenv.Load()
	apiKey := os.Getenv("BALLDONTLIE_API_KEY")
	if apiKey == "" {
		log.Fatal("API Key is missing! Set BALLDONTLIE_API_KEY in your environment.")
	}
	return apiKey
}

// Fetch JSON from API
func fetchAPI(apiURL string, target interface{}) error {
	apiKey := loadAPIKey()
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// HTTP Handler for searching players
func playerNameHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/players/")
	if id == "" {
		http.Error(w, "Missing player ID", http.StatusBadRequest)
		return
	}
	query := r.URL.Query().Get("search")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	apiURL := "https://api.balldontlie.io/v1/players?search=" + url.QueryEscape(query)
	var result PlayerResponse

	err := fetchAPI(apiURL, &result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching players: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HTTP Handler for searching players by ID
func playerIdHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/players/")
	if id == "" {
		http.Error(w, "Missing player ID", http.StatusBadRequest)
		return
	}

	apiURL := "https://api.balldontlie.io/v1/players/" + id
	var result PlayerIdResponse

	err := fetchAPI(apiURL, &result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching players: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/players", playerNameHandler) //Search by name
	mux.HandleFunc("/players/", playerIdHandler)  //Search by id

	// Enable CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow requests from React
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(mux)

	port := "8080"
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
