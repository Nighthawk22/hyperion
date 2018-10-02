package main

import (
	"os"
	"time"

	"github.com/Nighthawk22/hyperion/pkg/concourse"
	"github.com/Nighthawk22/hyperion/pkg/led"
	"github.com/Nighthawk22/hyperion/pkg/prometheus"
	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stdout)
	ledClient := led.New(led.Config{
		Normal:  led.RGB{255, 255, 255},
		Warning: led.RGB{255, 255, 255},
		Error:   led.RGB{255, 255, 255},
		Log:     log,
	})

	adaptor := ledClient.NewAdaptor()
	ledStripTop := ledClient.NewLEDStrip(adaptor, "13", "11", "15")
	ledStripBottom := ledClient.NewLEDStrip(adaptor, "29", "31", "33")

	concourseClient := concourse.New(concourse.Config{
		URL: "https://taa-ci-01.local.netconomy.net",
		Log: log,
	})

	prometheusClient := prometheus.New(prometheus.Config{
		URL: "https://alertmanager.monitoring.tools.local.netconomy.net",
		Log: log,
	})

	for true {
		prometheusClient.
			time.Sleep(1 * time.Minute)
	}

}

//https: //taa-ci-01.local.netconomy.net/api/v1/jobs
