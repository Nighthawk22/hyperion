// Copyright 2017 Hyperion Team
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License.  You may obtain a copy
// of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
// License for the specific language governing permissions and limitations under
// the License.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

var alertCounter int

const (
	redPort    = "17"
	greenPort  = "22"
	bluePort   = "24"
	lightBlue  = "50"
	lightRed   = "200"
	light      = "20"
	dark       = "0"
	ledCommand = "/usr/local/bin/pigs"
)

type alertManagerAlerts struct {
	Status string            `json:"status"`
	Alerts []prometheusAlert `json:"data"`
}

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
	log.Printf("Changing LED to alert!")
	err := exec.Command(ledCommand, "p", bluePort, dark).Run()
	if err != nil {
		log.Println(err)
	}

	err = exec.Command(ledCommand, "p", redPort, lightRed).Run()
	if err != nil {
		log.Println(err)
	}
}

func changeLEDToNormal() {
	log.Printf("Changing alert to normal!")
	err := exec.Command(ledCommand, "p", redPort, dark).Run()
	if err != nil {
		log.Println(err)
	}
	err = exec.Command(ledCommand, "p", bluePort, lightBlue).Run()
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

func startWebServer() {
	alertCounter = 0
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
	log.Println("Serving on Port 8080")
}

func requestAlertManager(alertManagerURL string) (alertManagerAlerts, error) {
	var alertManagerResponse alertManagerAlerts
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/alerts", alertManagerURL))
	if err != nil {
		log.Printf("Could not retrieve api")
		return alertManagerResponse, err
	}
	err = json.NewDecoder(resp.Body).Decode(&alertManagerResponse)
	return alertManagerResponse, err
}

func countAlertsFromManagerResponse(alertManagerReponse alertManagerAlerts) {
	var isAlert bool
	for _, alert := range alertManagerReponse.Alerts {
		if alert.EndsAt == "0001-01-01T00:00:00Z" {
			changeLEDToAlert()
			isAlert = true
			return
		}
	}

	if !isAlert {
		changeLEDToNormal()
	}
}

func pollingAlertManager(alertManagerURL string, pollingInterval int) {
	log.Printf("Start polling alertmanager on url %s", alertManagerURL)
	ticker := time.NewTicker(time.Duration(pollingInterval) * time.Second)

	for range ticker.C {
		alertManagerResponse, err := requestAlertManager(alertManagerURL)
		if err != nil {
			log.Print(err)
		}
		countAlertsFromManagerResponse(alertManagerResponse)
	}
}

func main() {
	changeLEDToNormal()
	webServerMod := flag.Bool("web", false, "Starting Webserver for pushing alerts")
	alertManagerURL := flag.String("url", "", "url of the Altermanager for polling alert status")
	pollingInterval := flag.Int("interval", 3, "Polling interval in seconds")
	flag.Parse()

	if *webServerMod {
		startWebServer()
	} else {
		if *alertManagerURL == "" {
			log.Fatal("Flags not provided. Aborting")
		} else {
			pollingAlertManager(*alertManagerURL, *pollingInterval)
		}
	}
}
