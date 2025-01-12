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
	Status       string `json:"status"`
  StatusAdd       string `json:"status_add_user"`
}
func reloadConfig() {
	url := "https://serverproxy.vncloud.net/vnproxyrandom.php?action=config"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}
	result := string(body)
	if result == "error" {
		log.Println("Operation failed: Server returned 'error'")
		return
	}

	filePath := "/etc/squid/squid.conf"
	err = ioutil.WriteFile(filePath, []byte(result), 0644)
	if err != nil {
		log.Fatalf("Failed to write to file %s: %v", filePath, err)
	}
	log.Printf("Successfully updated %s", filePath)
	cmd := exec.Command("systemctl", "restart", "squid")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to restart squid: %v", err)
	}
	log.Println("Successfully restarted squid")
}

func processUsers() {

	resp, err := http.Get("https://serverproxy.vncloud.net/vnproxyrandom.php?action=cron")
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
		if user.StatusAdd == "0" {
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

				createURL := fmt.Sprintf("https://serverproxy.vncloud.net/vnproxyrandom.php?adduser=%s", user.ID)
				respCreate, err := http.Get(createURL)
				if err != nil {
					log.Printf("Error Request Update for User %s: %v", user.Username, err)
				}
				respCreate.Body.Close()
			}
		}

		if user.Status == "expired" {
			if user.StatusAdd == "1" {
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
				deleteURL := fmt.Sprintf("https://serverproxy.vncloud.net/user.php?deleteuser=%s", user.ID)
				respDelete, err := http.Get(deleteURL)
				if err != nil {
					log.Printf("Error Request Update for User %s: %v", user.Username, err)
				} else {
					log.Printf("Restarting Squid Success")
				}
				respDelete.Body.Close()
			}
		}
	  if user.Status == "suspend" {
			if user.StatusAdd == "1" {
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
				deleteURL := fmt.Sprintf("https://serverproxy.vncloud.net/user.php?deleteuser=%s", user.ID)
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
    reloadConfig()
		time.Sleep(1 * time.Second)
	}
}

