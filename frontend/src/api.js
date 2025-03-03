const API_URL = "http://localhost:8080";

export async function searchPlayerByName(name) {
  const response = await fetch(`${API_URL}/players?search=${name}`);
  if (!response.ok) {
    throw new Error("Failed to fetch players");
  }
  return response.json();
}


export async function searchPlayerById(id) {
    const response = await fetch(`${API_URL}/players/${id}`); // Ensure the ID is appended properly
    if (!response.ok) {
      throw new Error(`Failed to fetch player by ID: ${response.statusText}`);
    }
    return response.json();
  }
  
