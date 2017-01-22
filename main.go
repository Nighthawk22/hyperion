package main

import (
	"encoding/json"
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
	light      = "20"
	dark       = "0"
	ledCommand = "/usr/local/bin/pigs"
)

type alertManagerAlerts struct {
	Status 		string `json:"status"`
	Alerts		[]prometheusAlert	`json:"data"`
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

func startWebServer() {
	alertCounter = 0
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
	log.Println("Serving on Port 8080")
}

func requestAlertManager(alertManagerIp string) {
	resp, err := http.Get(fmt.Sprint("%s/api/v1/alerts"))
	if err != nil {
		log.Printf("Could not retrieve api")
		return
	}
	var alertManagerResponse []alertManagerAlerts
	json.NewDecoder(resp.Body).Decode(alertManagerResponse)


}


func pollingAlertManager(alertManagerIp string, pollingInterval int) {
	ticker := time.NewTicker(pollingInterval * time.Second)
	go func() {
		for t := range ticker.C {
			requestAlertManager(alertManagerIp)
		}
	}
}


func main() {
	changeLEDToNormal()
	webServerMod := flag.Bool("web", false, "Starting Webserver for pushing alerts")
	alertManagerIp := flag.String("url", "", "url of the Altermanager for polling alert status")
	pollingInterval := flag.Int("interval", 3, "Polling interval in seconds")

	if webServerMod {
		startWebServer()
	} else {
		if !alertManagerIp || !pollingInterval {
			log.Fatal("Flags not provided. Aborting")
		} else {
			pollingAlertManager(*alertManagerIp, *pollingInterval)
		}
	}


	
}
