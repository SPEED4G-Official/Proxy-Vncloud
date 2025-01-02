package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Expired      string `json:"expired"`
	Status       string `json:"status"`
	Quantity     string `json:"quantity"`
	Phone        string `json:"phone"`
	CreateServer string `json:"create_server"`
	Array        string `json:"array"`
}

func processUsers() {
	resp, err := http.Get("https://httpproxy.vncloud.net/user.php?list")
	if err != nil {
		log.Fatalf("Error Request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error Reading Body: %v", err)
	}
	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		log.Fatalf("Error Load JSON: %v", err)
	}

	for _, user := range users {
		if user.CreateServer == "0" {
			if user.Status == "active" {
				fmt.Printf("Create User: %s\n", user.Username)
				cmd := exec.Command("sudo", "/usr/bin/htpasswd", "-b", "/etc/squid/passwd", user.Username, user.Password)
				err := cmd.Run()
				if err != nil {
					log.Printf("Error Adding User %s: %v", user.Username, err)
				}

				cmdRestart := exec.Command("sudo", "systemctl", "restart", "squid")
				err = cmdRestart.Run()
				if err != nil {
					log.Printf("Error Restarting Squid: %v", err)
				} else {
					log.Printf("Restarting Squid Success")
				}

				createURL := fmt.Sprintf("https://httpproxy.vncloud.net/user.php?createproxy=%s", user.ID)
				respCreate, err := http.Get(createURL)
				if err != nil {
					log.Printf("Error Request Update for User %s: %v", user.Username, err)
				}
				respCreate.Body.Close()
			}
		}

		if user.Status == "expired" {
			if user.CreateServer == "1" {
				fmt.Printf("Delete User: %s\n", user.Username)
				cmd := exec.Command("sudo", "/usr/bin/htpasswd", "-D", "/etc/squid/passwd", user.Username)
				err := cmd.Run()
				if err != nil {
					log.Printf("Error Deleting User %s: %v", user.Username, err)
				}

				cmdRestart := exec.Command("sudo", "systemctl", "restart", "squid")
				err = cmdRestart.Run()
				if err != nil {
					log.Printf("Error Restarting Squid: %v", err)
				}
				deleteURL := fmt.Sprintf("https://httpproxy.vncloud.net/user.php?deleteproxy=%s", user.ID)
				respDelete, err := http.Get(deleteURL)
				if err != nil {
					log.Printf("Error Request Update for User %s: %v", user.Username, err)
				} else {
					log.Printf("Restarting Squid Success")
				}
				respDelete.Body.Close()
			}
		}
	}
}

func main() {
	for {
		processUsers()
		time.Sleep(30 * time.Second)
	}
}
