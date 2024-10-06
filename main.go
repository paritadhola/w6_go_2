package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Create the IIT Server Management System Data Structure
type Server struct {
	ID        int    `json:"ID"`
	HostName  string `json:"HostName"`
	IPAddress string `json:"IPAddress"`
	Status    string `json:"Status"`
}

var serverList []Server
var nextServerID int = 1

// Input method for Server Data
func inputServerData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var server Server
	// Read body of the POST request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid read request body", http.StatusBadRequest)
		return
	}
	// Parse the JSON data
	server.ID = nextServerID
	nextServerID++
	err = json.Unmarshal(body, &server)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	serverList = append(serverList, server)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(server)
}

// get all server data from server system
func getAllServerData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serverList)
}

func getServerByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract server ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/getServerByID/")
	serverID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}

	// Search for the server by ID
	for _, server := range serverList {
		if server.ID == serverID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(server)
			return
		}
	}

	// If server not found, return 404
	http.Error(w, "Server not found", http.StatusNotFound)
}
func updateServerByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract server ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/updateServerByID/")
	serverID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}

	// Find the server by ID
	for i, server := range serverList {
		if server.ID == serverID {
			// Parse the updated server data
			var updatedServer Server
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Invalid read request body", http.StatusBadRequest)
				return
			}

			err = json.Unmarshal(body, &updatedServer)
			if err != nil {
				http.Error(w, "Invalid JSON format", http.StatusBadRequest)
				return
			}

			// Update server data
			serverList[i].HostName = updatedServer.HostName
			serverList[i].IPAddress = updatedServer.IPAddress
			serverList[i].Status = updatedServer.Status

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(serverList[i])
			return
		}
	}
}
func deleteServerByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract server ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/deleteServerByID/")
	serverID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}

	// Search for the server by ID and delete it
	for i, server := range serverList {
		if server.ID == serverID {
			// Delete the server from the list
			serverList = append(serverList[:i], serverList[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Server deleted successfully"})
			return
		}
	}

	// If server not found, return 404
	http.Error(w, "Server not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/getAllServerData", getAllServerData)
	http.HandleFunc("/inputServerData", inputServerData)
	http.HandleFunc("/getServerByID/", getServerByID)
	http.HandleFunc("/updateServerByID/", updateServerByID)
	http.HandleFunc("/deleteServerByID/", deleteServerByID)
	fmt.Println("Server running on port : 46664")
	log.Fatal(http.ListenAndServe(":46664", nil))
}
