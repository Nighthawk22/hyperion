# Hyperion: An LED Alertmanager

Hyperion is used for polling alerts form the Prometheus
Alertmanager and getting build status of the concourse ci server.

The idea is to use the dynatrace ufo https://github.com/Dynatrace/ufo-esp32 with a rasperry pi as controller and two led stripes.
To get led stripes working with the raspberry pi we used this manual https://dordnung.de/raspberrypi-ledstrip/ 
and added an additional stripe.

## Prometheus alertmanager

If an alert is raised the color of the top LED strip changes to error. When there is no alert left it changes back to normal.

## Concourse ci

On an running job the led at the bottom changes to a pulsating warning. When a job has an error state and there 
is no running job it changes to error. If there is no error job and also no running job the leds change to normal.


## Usage

Build your config file based on `config.example.yaml` and start hyperion with the `--config config.yaml` flag.
