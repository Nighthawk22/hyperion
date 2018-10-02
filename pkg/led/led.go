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
	"context"
	"time"

	"github.com/rs/zerolog"
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

func (c Client) NewLEDStrip(adaptor gpio.DigitalWriter, name string, redPin string, greenPin string, bluePin string) *LEDStrip {
	return &LEDStrip{
		LedDriver: gpio.NewRgbLedDriver(adaptor, redPin, greenPin, bluePin),
		Name:      name,
	}
}

func (c Client) On(led *LEDStrip) error {
	c.config.Log.Info().Str("led", led.Name).Msg("Setting to on")
	return led.LedDriver.On()
}

func (c Client) Normal(led *LEDStrip) error {
	c.config.Log.Info().Str("led", led.Name).Msg("Setting to normal")
	return led.LedDriver.SetRGB(c.config.Normal.Red, c.config.Normal.Green, c.config.Normal.Blue)
}

func (c Client) Warning(led *LEDStrip) error {
	c.config.Log.Info().Str("led", led.Name).Msg("Setting to warning")
	return led.LedDriver.SetRGB(c.config.Warning.Red, c.config.Warning.Green, c.config.Warning.Blue)
}

func (c Client) Error(led *LEDStrip) error {
	c.config.Log.Info().Str("led", led.Name).Msg("Setting to error")
	return led.LedDriver.SetRGB(c.config.Error.Red, c.config.Error.Green, c.config.Error.Blue)
}

//TODO: Use cancelcontext
func (c Client) Pulsating(ctx context.Context, led *LEDStrip, rgb RGB) {
	led.LedDriver.SetRGB(rgb.Red, rgb.Green, rgb.Blue)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			pulse(led)
		}

	}
}

func pulse(led *LEDStrip) {
	time.Sleep(2 * time.Second)
	led.LedDriver.Toggle()
}
