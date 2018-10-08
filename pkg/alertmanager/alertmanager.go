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
package alertmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
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

//Config is the concourse base config
type Config struct {
	URL string
	Log zerolog.Logger
}

//Client represents the concourse client.
type Client struct {
	config *Config
}

//New creates a new concourse client
func New(config Config) *Client {
	return &Client{
		config: &config,
	}
}

//CheckAlerts returns if there is an alert currently.
func (c Client) CheckAlerts(ctx context.Context) (bool, error) {
	var alertManagerResponse alertManagerAlerts
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/alerts", c.config.URL))
	if err != nil {
		return false, err
	}
	err = json.NewDecoder(resp.Body).Decode(&alertManagerResponse)
	for _, alert := range alertManagerResponse.Alerts {
		if alert.EndsAt == "0001-01-01T00:00:00Z" {
			return true, nil
		}
	}

	return false, nil
}
