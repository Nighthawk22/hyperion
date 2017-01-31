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
	"time"

	"github.com/nighthawk22/hyperion/led"
)

const (
	redPort   = "17"
	greenPort = "22"
	bluePort  = "24"
	lightBlue = "50"
	lightRed  = "200"
	light     = "20"
	dark      = "0"
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

func changeLEDToAlert() {
	log.Printf("Changing LED to alert!")
	if err := led.ChangeLED(bluePort, dark); err != nil {
		log.Println(err)
	}

	if err := led.ChangeLED(redPort, lightRed); err != nil {
		log.Println(err)
	}
}

func changeLEDToNormal() {
	log.Printf("Changing alert to normal!")
	if err := led.ChangeLED(redPort, dark); err != nil {
		log.Println(err)
	}

	if err := led.ChangeLED(bluePort, lightBlue); err != nil {
		log.Println(err)
	}
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

	alertManagerURL := flag.String("url", "", "url of the Altermanager for polling alert status")
	pollingInterval := flag.Int("interval", 3, "Polling interval in seconds")
	flag.Parse()

	if *alertManagerURL == "" {
		log.Fatal("Flags not provided. Aborting")
	} else {
		pollingAlertManager(*alertManagerURL, *pollingInterval)
	}

}
