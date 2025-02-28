package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
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

type TeamResponse struct {
	Data []Team `json:"data"`
}

type PlayerResponse struct {
	Data []Player `json:"data"`
}

type PlayerIdResponse struct {
	//single player object, not an array
	Data Player `json:"data"`
}

// Load API Key from .env
func loadAPIKey() string {
	// Load .env file (if exists)
	godotenv.Load()

	apiKey := os.Getenv("BALLDONTLIE_API_KEY")
	if apiKey == "" {
		log.Fatal("API Key is missing! Set BALLDONTLIE_API_KEY in your environment.")
	}

	return apiKey
}

// Fetch JSON from API with Authorization Header
func fetchAPI(apiURL string, target interface{}) error {
	apiKey := loadAPIKey()
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}

	// Set API Key in Headers
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

// Search for teams
func searchTeam(teamName string) {
	apiURL := "https://api.balldontlie.io/v1/teams"
	var result TeamResponse

	err := fetchAPI(apiURL, &result)
	if err != nil {
		log.Fatalf("Error fetching teams: %v", err)
	}

	// Search for the team
	found := false
	for _, team := range result.Data {
		if strings.Contains(strings.ToLower(team.Name), strings.ToLower(teamName)) {
			fmt.Printf("Found Team: %d: %s\n", team.ID, team.Name)
			found = true
		}
	}

	if !found {
		fmt.Println("Team not found.")
	}
}

// Search for players
func searchPlayer(playerName string) {
	apiURL := "https://api.balldontlie.io/v1/players?search=" + url.QueryEscape(playerName)
	var result PlayerResponse

	err := fetchAPI(apiURL, &result)
	if err != nil {
		log.Fatalf("Error fetching players: %v", err)
	}

	// Display players found
	if len(result.Data) == 0 {
		fmt.Println("No players found.")
		return
	}

	for _, player := range result.Data {
		fmt.Printf("ID: %d | Player: %s %s | Team: %s | Stats: %s %s #%s\n", player.ID, player.First, player.Last, player.Team.Name, player.Position, player.Height, player.Jersey)
	}
}

// Search for players by ID
func searchPlayerId(playerId string) {
	apiURL := "https://api.balldontlie.io/v1/players/" + url.QueryEscape(playerId)
	var player PlayerIdResponse

	err := fetchAPI(apiURL, &player)
	if err != nil {
		log.Fatalf("Error fetching players: %v", err)
	}

	// Display player details
	fmt.Printf("ID: %d | Player: %s %s | Team: %s | Stats: %s %s #%s\n",
		player.Data.ID, player.Data.First, player.Data.Last, player.Data.Team.Name, player.Data.Position, player.Data.Height, player.Data.Jersey)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1: Search for a Team")
		fmt.Println("2: Search for a Player")
		fmt.Println("3: Search for Player by ID")
		fmt.Println("4: Exit")

		fmt.Print("Enter choice: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			fmt.Print("Enter team name: ")
			scanner.Scan()
			teamName := scanner.Text()
			searchTeam(teamName)

		case "2":
			fmt.Print("Enter player name: ")
			scanner.Scan()
			playerName := scanner.Text()
			searchPlayer(playerName)

		case "3":
			fmt.Print("Enter player ID: ")
			scanner.Scan()
			playerId := scanner.Text()
			searchPlayerId(playerId)

		case "4":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice, try again.")
		}
	}
}
