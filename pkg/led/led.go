package led

import (
	"time"

	"github.com/rs/zerolog"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

type RGB struct {
	Red   byte
	Green byte
	Blue  byte
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
}

func New(c Config) *Client {
	return &Client{
		&c,
	}
}

func (c Client) NewAdaptor() *raspi.Adaptor {
	return raspi.NewAdaptor()
}

func (c Client) NewLEDStrip(adaptor gpio.DigitalWriter, redPin string, greenPin string, bluePin string) *LEDStrip {
	return &LEDStrip{
		LedDriver: gpio.NewRgbLedDriver(adaptor, redPin, greenPin, bluePin),
	}
}

func (c Client) On(led *LEDStrip) error {
	return led.LedDriver.On()
}

func (c Client) Normal(led *LEDStrip) error {
	return led.LedDriver.SetRGB(c.config.Normal.Red, c.config.Normal.Green, c.config.Normal.Blue)
}

func (c Client) Warning(led *LEDStrip) error {
	return led.LedDriver.SetRGB(c.config.Warning.Red, c.config.Warning.Green, c.config.Warning.Blue)
}

func (c Client) Error(led *LEDStrip) error {
	return led.LedDriver.SetRGB(c.config.Error.Red, c.config.Error.Green, c.config.Error.Blue)
}

//TODO: Use cancelcontext
func (c Client) Pulsating(led *LEDStrip, rgb RGB) func() int {
	return func() int {
		led.LedDriver.SetRGB(rgb.Red, rgb.Green, rgb.Blue)
		for true {
			time.Sleep(2 * time.Second)
			led.LedDriver.Toggle()
		}
		return 1
	}
}
