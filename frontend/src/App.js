import { useState } from "react";
import { searchPlayerByName } from "./api";
import "./App.css";

function App() {
  const [input, setInput] = useState("");
  const [players, setPlayers] = useState([]);

  const handleSearch = async () => {
    if (!input.trim()) return; // Prevent empty searches

    try {
      const data = await searchPlayerByName(input);
      setPlayers(data.data); // Assuming the API response contains a "data" array
    } catch (error) {
      console.error("Error fetching player data:", error);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>NBA Player Search</h1>
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Enter player name"
        />
        <button onClick={handleSearch}>Search</button>

        <div>
          {players.length > 0 ? (
            players.map((player) => (
              <div key={player.id} className="player-card">
                <h3>{player.first_name} {player.last_name}</h3>
                <p>Team: {player.team?.full_name || "N/A"}</p>
                <p>Position: {player.position || "N/A"}</p>
              </div>
            ))
          ) : (
            <p>No players found.</p>
          )}
        </div>
      </header>
    </div>
  );
}

export default App;
