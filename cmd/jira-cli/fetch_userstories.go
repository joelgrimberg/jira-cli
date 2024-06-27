package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func fetchStories() {
	url := "https://api.clubhouse.io/api/v3/userstories"
	f, err := os.ReadFile("config.json")
	if err != nil {
		log.Println(err)
	}

	var data map[string]interface{}
	json.Unmarshal([]byte(f), &data)

	log.Println(data)
	for k, v := range data {
		log.Println(k, ":", v)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Add any necessary headers to authenticate or provide authorization

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println(string(body))
}
