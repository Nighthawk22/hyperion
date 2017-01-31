# Hyperion: An LED Alertmanager

Hyperion is used for polling or receiving alerts form the Prometheus
Alertmanager.  If an alert is found then the color of the LEDs changes to
red. If there is no alert the color will be blue.

Right now two operational modes are supported:

## Receiving

For the push mode add the `--web` flag and a webserver on port 8080 will be
started.  The webserver waits for payloads similar
to
[Prometheus webhooks](https://prometheus.io/docs/alerting/configuration/#webhook-receiver-%3Cwebhook_config%3E) send
from the Alertmanager. The Alertmanager also has to send resolve requests.

## Polling

If you just provide a `--url` flag hyperion will poll the Prometheus
Alertmanager every 3 seconds for new alerts.

## Flags

 * `--web`: Starts a webserver waiting for new alerts.
 * `--url`: The url for polling the prometheus
   alertmanager. E.g. https://alertmanager.example.com
 * `--interval`: Polling interval. Defaults to every 3 seconds.
