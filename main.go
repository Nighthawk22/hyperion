package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

var alertCounter int

const (
	redPort    = "17"
	greenPort  = "22"
	bluePort   = "24"
	light      = "20"
	dark       = "0"
	ledCommand = "/usr/local/bin/pigs"
)

type prometheusAlert struct {
	Labels      string `json:"labels"`
	Annotations string `json:"annotations"`
	StartsAt    string `json:"startsAt"`
	EndsAt      string `json:"endsAt"`
	Status      string `json:"status"`
}

type prometheusAlertNotification struct {
	GroupKey string            `json:"groupKey"`
	Receiver string            `json:"receiver"`
	Status   string            `json:"status"`
	Alerts   []prometheusAlert `json:"alerts"`
}

func countAlerts(alertNotification prometheusAlertNotification) {
	for _, alert := range alertNotification.Alerts {
		if alert.Status == "resolved" {
			alertCounter--
			if alertCounter < 0 {
				alertCounter = 0
			}
		} else {
			alertCounter++
		}
	}
}

func changeLEDToAlert() {
	err := exec.Command(ledCommand, "p", bluePort, dark).Run()
	if err != nil {
		log.Println(err)
	}

	err = exec.Command(ledCommand, "p", redPort, light).Run()
	if err != nil {
		log.Println(err)
	}
}

func changeLEDToNormal() {
	err := exec.Command(ledCommand, "p", redPort, dark).Run()
	if err != nil {
		log.Println(err)
	}
	err = exec.Command(ledCommand, "p", bluePort, light).Run()
	if err != nil {
		log.Println(err)
	}
}

func operateLED() {
	if alertCounter > 0 {
		changeLEDToAlert()
	} else {
		changeLEDToNormal()
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var alertNotification prometheusAlertNotification
	err := decoder.Decode(&alertNotification)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Bad Request")
	}
	defer r.Body.Close()
	countAlerts(alertNotification)
	operateLED()

	fmt.Fprintf(w, "NR of alerts "+strconv.Itoa(alertCounter))

}
func main() {
	alertCounter = 0
	changeLEDToNormal()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
	log.Println("Serving on Port 8080")
}
