# LED Alertmanager

LED Alertmanager is used for polling or receiving alerts form the prometheus alertmanager.  
If an alert is found the color of the leds changes to red. If there is no alert the color will be blue.

Two methods were added:

## Receiving

For the Pushing method add the `--web` flag and a webserver on port 8080 will be startet. 
The webserver waits for webhooks (https://prometheus.io/docs/alerting/configuration/#webhook-receiver-%3Cwebhook_config%3E)
send from the alertmanager. The alertmanager also has to send resolbe requests.

## Polling

If you just provide a `--url` flag the led-alertmanager will poll the prometheus alertmanager every 3 seconds for new alerts.

## Flags

 * `--web` Starts a webserver waiting for new alerts.
 * `--url` The url for polling the prometheus alertmanager. E.g. https://alertmanager.example.com
 * `--interval` Polling interval. Defaults to every 3 seconds.

## TODO

 * Make colors configurable.
 * Better method for alert handling in the receiving mode. 