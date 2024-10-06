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

// Input the server data
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

// Get server by ID
func getServerByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the server ID from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Server ID not provided", http.StatusBadRequest)
		return
	}

	// Convert the server ID from string to integer
	serverID, err := strconv.Atoi(pathParts[2])
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
	//	http.Error(w, "Server not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/getAllServerData", getAllServerData)
	http.HandleFunc("/inputServerData", inputServerData)
	http.HandleFunc("/getServerByID/", getServerByID)
	fmt.Println("Server running on port : 46664")
	log.Fatal(http.ListenAndServe(":46664", nil))
}
