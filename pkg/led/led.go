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
package led

import (
	"time"

	"github.com/rs/zerolog"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

type RGB struct {
	Red   byte `yaml:"red"`
	Green byte `yaml:"green"`
	Blue  byte `yaml:"blue"`
}

type Config struct {
	Normal  RGB
	Warning RGB
	Error   RGB
	Log     zerolog.Logger
}

type Client struct {
	config *Config
}

type LEDStrip struct {
	LedDriver *gpio.RgbLedDriver
	Adaptor   *raspi.Adaptor
	Name      string
}

func New(c Config) *Client {
	return &Client{
		&c,
	}
}

func (c Client) NewAdaptor() *raspi.Adaptor {
	return raspi.NewAdaptor()
}

func (c Client) NewLEDStrip(adaptor *raspi.Adaptor, name string, redPin string, greenPin string, bluePin string) *LEDStrip {
	ledDriver := gpio.NewRgbLedDriver(adaptor, redPin, greenPin, bluePin)
	ledDriver.SetName(name)
	return &LEDStrip{
		LedDriver: ledDriver,
		Adaptor:   adaptor,
		Name:      name,
	}
}

func (c Client) Normal(led *LEDStrip) error {
	c.config.Log.Info().Str("led", led.LedDriver.Name()).Msg("Setting to normal")
	return led.LedDriver.SetRGB(c.config.Normal.Red, c.config.Normal.Green, c.config.Normal.Blue)
}

func (c Client) Warning(led *LEDStrip) error {
	c.config.Log.Info().Str("led", led.LedDriver.Name()).Msg("Setting to warning")
	return led.LedDriver.SetRGB(c.config.Warning.Red, c.config.Warning.Green, c.config.Warning.Blue)
}

func (c Client) Error(led *LEDStrip) error {
	c.config.Log.Info().Str("led", led.LedDriver.Name()).Msg("Setting to error")
	return led.LedDriver.SetRGB(c.config.Error.Red, c.config.Error.Green, c.config.Error.Blue)
}

func (c Client) Pulsating(led *LEDStrip, rgb RGB) *time.Ticker {
	c.config.Log.Info().Str("led", led.LedDriver.Name()).Msg("Setting to pulsating")
	led.LedDriver.SetRGB(rgb.Red, rgb.Green, rgb.Blue)
	return gobot.Every(1*time.Second, func() {
		c.config.Log.Info().Msg("Toggle")
		led.LedDriver.Toggle()
	})

}
