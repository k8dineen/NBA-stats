import { useState } from "react";
import { searchPlayerByName } from "./api";
import { searchPlayerById } from "./api";
import "./App.css";

function App() {
  const [input, setInput] = useState("");
  const [searchType, setSearchType] = useState("name"); // "name" or "id"
  const [players, setPlayers] = useState([]);

  const handleSearch = async () => {
    if (!input.trim()) return; // Prevent empty searches

    try {
      if (searchType === "name") {
        const data = await searchPlayerByName(input);
        setPlayers(data.data || []); 
      } else {
        const id = input.trim(); // Remove spaces to avoid issues
        const data = await searchPlayerById(id);
        setPlayers([data.data]);
      }
    } catch (error) {
      console.error("Error fetching player data:", error);
      setPlayers([]);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h6>NBA-stats</h6>
      
        {/* Radio Buttons for Selecting Search Type */}
        <div>
          <label>
            <input
              type="radio"
              value="name"
              checked={searchType === "name"}
              onChange={() => setSearchType("name")}
            />
            Search by Name
          </label>
          <label>
            <input
              type="radio"
              value="id"
              checked={searchType === "id"}
              onChange={() => setSearchType("id")}
            />
            Search by ID
          </label>
          
          {/* Search Input & Button */}
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder={searchType === "name" ? "Enter player name" : "Enter player ID"}
          />
          <button onClick={handleSearch}>Search</button>
        </div>
      </header>

      {/* Player Results */}
      <div>
        {players.length > 0 ? (
          players.map((player) => (
            <div key={player.id} className="player-card">
              <h3>{player.first_name} {player.last_name}</h3>
              <p><strong>ID:</strong> {player.id}</p>
              <p><strong>Team:</strong> {player.team?.full_name || "N/A"}</p>
              <p><strong>Position:</strong> {player.position || "N/A"}</p>
              <p><strong>Height:</strong> {player.height || "N/A"}</p>
              <p><strong>Weight:</strong> {player.weight + "lbs." || "N/A"}</p>
            </div>
          ))
        ) : (
          <p>No players found.</p>
        )}
      </div>
    </div>
  );
}

export default App;
