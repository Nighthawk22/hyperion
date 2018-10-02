// Copyright 2018 Hyperion Team
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
	"context"
	"os"
	"time"

	"github.com/Nighthawk22/hyperion/pkg/alertmanager"
	"github.com/Nighthawk22/hyperion/pkg/concourse"
	"github.com/Nighthawk22/hyperion/pkg/led"
	"github.com/rs/zerolog"
)

type LedStripConfig struct {
	Name     string `yaml:"string"`
	RedPin   string `yaml:"redPin"`
	GreenPin string `yaml:"greenPin"`
	BluePin  string `yaml:"bluePin"`
}
type Config struct {
	Alertmanager struct {
		URL string `yaml:"url"`
	} `yaml:"alertmanager"`
	Concourse struct {
		URL string `yaml:"url"`
	} `yaml:"concourse"`
	LedStrips struct {
		Top    LedStripConfig `yaml:"top"`
		Bottom LedStripConfig `yaml:"bottom"`
	} `yaml:"ledStrips"`
	Colors struct {
		Normal  led.RGB `yaml:"normal"`
		Warning led.RGB `yaml:"warning"`
		Error   led.RGB `yaml:"Error"`
	}
}

//Run starts a new hyperion instance
func Run(config Config) {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	log.Info().Interface("Config", config).Msg("Config loaded")

	ledClient := led.New(led.Config{
		Normal:  config.Colors.Normal,
		Warning: config.Colors.Warning,
		Error:   config.Colors.Error,
		Log:     log,
	})

	adaptor := ledClient.NewAdaptor()
	ledStripTop := ledClient.NewLEDStrip(adaptor, config.LedStrips.Top.Name, config.LedStrips.Top.RedPin, config.LedStrips.Top.GreenPin, config.LedStrips.Top.BluePin)
	ledStripBottom := ledClient.NewLEDStrip(adaptor, config.LedStrips.Bottom.Name, config.LedStrips.Bottom.RedPin, config.LedStrips.Bottom.GreenPin, config.LedStrips.Bottom.BluePin)

	concourseClient := concourse.New(concourse.Config{
		URL: config.Concourse.URL,
		Log: log,
	})

	alertManagerClient := alertmanager.New(alertmanager.Config{
		URL: config.Alertmanager.URL,
		Log: log,
	})

	blinkingContext, cancel := context.WithCancel(context.Background())

	for {
		log.Info().Msg("Checking concourse")
		running, errJobs, err := concourseClient.CheckJobs(context.Background())

		if err != nil {
			log.Error().Err(err).Msg("Could not call concourse")
			err = ledClient.Error(ledStripBottom)
		} else {
			if running {
				ledClient.Pulsating(blinkingContext, ledStripBottom, config.Colors.Warning)
			} else if errJobs {
				cancel()
				err = ledClient.Error(ledStripBottom)
			} else {
				cancel()
				err = ledClient.Normal(ledStripBottom)
			}
		}

		if err != nil {
			log.Error().Err(err).Msg("Could not set led strip.")
		}

		log.Info().Msg("Checking prometheus")

		promAlert, err := alertManagerClient.CheckAlerts(context.Background())

		if err != nil {
			log.Error().Err(err).Msg("Could not call prometheus")
			err = ledClient.Error(ledStripTop)
		} else {
			if promAlert {
				err = ledClient.Error(ledStripTop)
			} else {
				err = ledClient.Normal(ledStripTop)
			}
		}

		if err != nil {
			log.Error().Err(err).Msg("Could not set led strip.")
		}

		time.Sleep(1 * time.Minute)
	}
}
