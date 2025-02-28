package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"full_name"`
}

type APIResponse struct {
	Data []Team `json:"data"`
}

func getNBATeams() {
	// Load .env file (if exists)
	godotenv.Load()

	apiKey := os.Getenv("BALLDONTLIE_API_KEY")
	if apiKey == "" {
		log.Fatal("API Key is missing! Set BALLDONTLIE_API_KEY in your environment.")
	}

	url := "https://api.balldontlie.io/v1/teams"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Corrected: Use only API key (no "Bearer")
	req.Header.Set("Authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API request failed with status: %s", resp.Status)
	}

	var result APIResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Error decoding response: %v", err)
	}

	fmt.Println("NBA Teams:")
	for _, team := range result.Data {
		fmt.Printf("- %s\n", team.Name)
	}
}

func main() {
	getNBATeams()
}
