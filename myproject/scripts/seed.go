//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	base := "http://localhost:8080/api/v1"

	users := []map[string]string{
		{"username": "alice", "email": "alice@example.com", "password": "password123"},
		{"username": "bob", "email": "bob@example.com", "password": "password123"},
		{"username": "charlie", "email": "charlie@example.com", "password": "password123"},
	}

	for _, u := range users {
		data, _ := json.Marshal(u)
		resp, err := http.Post(base+"/register", "application/json", bytes.NewReader(data))
		if err != nil {
			fmt.Printf("Error registering %s: %v\n", u["username"], err)
			continue
		}
		fmt.Printf("Register %s: %d\n", u["username"], resp.StatusCode)
		resp.Body.Close()
	}

	fmt.Println("Seed complete.")
}