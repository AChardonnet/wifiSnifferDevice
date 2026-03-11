package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Device struct {
	Type string `json:"kismet.device.base.type"`
}

type Data struct {
	NbPeople int    `json:"nb_people"`
	Source   string `json:"source"`
}

func countWifiClients(host string, username string, password string, periodSec int) (int, error) {
	timestamp := time.Now().Unix() - int64(periodSec)
	url := fmt.Sprintf("http://%s/devices/views/phy-IEEE802.11/last-time/%d/devices.json", host, timestamp)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.SetBasicAuth(username, password)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Kismet API returned status: %s", resp.Status)
	}

	var devices []Device
	err = json.NewDecoder(resp.Body).Decode(&devices)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, d := range devices {
		if d.Type == "Wi-Fi Client" {
			count++
		}
	}

	return count, nil
}

func runEachSec(interval time.Duration, function func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		function()
	}
}

func runCounter() {
	t1 := time.Now()
	kismetUrl := os.Getenv("KISMET_URL")
	kismetUser := os.Getenv("KISMET_USER")
	kismetPassword := os.Getenv("KISMET_PASSWORD")
	count, err := countWifiClients(kismetUrl, kismetUser, kismetPassword, 600)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Wi-Fi Clients seen in the last 10 minutes:", count)

	data := Data{
		NbPeople: count,
		Source:   "wifi",
	}
	sendDataToBack(data)
	fmt.Printf("Execution took %d sec\n", int64(time.Since(t1).Seconds()))
}

func sendDataToBack(data Data) {
	backendUrl := os.Getenv("BACKEND_URL")
	url := fmt.Sprintf("http://%s/new_data", backendUrl)

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error :", err)
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)
}

func main() {
	fmt.Println("Loading .env ...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Counter starting")
	runEachSec(30*time.Second, runCounter)
}
