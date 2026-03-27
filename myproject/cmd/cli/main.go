package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"myproject/pkg/httpclient"
)

func main() {
	baseURL := flag.String("url", "http://localhost:8080", "Server base URL")
	action := flag.String("action", "health", "Action: health, register, login")
	email := flag.String("email", "", "User email")
	username := flag.String("username", "", "Username")
	password := flag.String("password", "", "Password")
	flag.Parse()

	client := httpclient.New(*baseURL)

	switch *action {
	case "health":
		body, status, err := client.Get("/healthz")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("[%d] %s\n", status, string(body))

	case "register":
		payload := map[string]string{
			"username": *username,
			"email":    *email,
			"password": *password,
		}
		data, _ := json.Marshal(payload)
		body, status, err := client.Post("/api/v1/register", data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("[%d] %s\n", status, string(body))

	case "login":
		payload := map[string]string{
			"email":    *email,
			"password": *password,
		}
		data, _ := json.Marshal(payload)
		body, status, err := client.Post("/api/v1/login", data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("[%d] %s\n", status, string(body))

	default:
		fmt.Fprintf(os.Stderr, "Unknown action: %s\n", *action)
		os.Exit(1)
	}
}