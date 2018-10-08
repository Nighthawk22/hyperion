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
package concourse

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/concourse/atc"
	"github.com/rs/zerolog"
)

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

//CheckJobs returns as first param if there are running jobs,
//as second if there are error jobs and as last if there was an error calling the api
func (c Client) CheckJobs(ctx context.Context) (bool, bool, error) {
	jobs, err := getJobs(ctx, c.config.URL)

	if err != nil {
		return false, false, err
	}

	c.config.Log.Info().Int("number", len(jobs)).Msg("Jobs received")

	var runningJobs bool

	if len(filterRunningJobs(jobs)) > 0 {
		runningJobs = true
	}

	var erroredJobs bool
	if len(filterErrorJobs(jobs)) > 0 {
		erroredJobs = true
	}

	return runningJobs, erroredJobs, nil
}

func getJobs(ctx context.Context, baseURL string) ([]atc.Job, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/jobs", baseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var jobs []atc.Job

	err = json.NewDecoder(res.Body).Decode(&jobs)

	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func filterErrorJobs(jobs []atc.Job) map[int]int {
	errorJobs := make(map[int]int)

	for _, job := range jobs {
		if job.FinishedBuild != nil {
			switch atc.BuildStatus(job.FinishedBuild.Status) {
			case atc.StatusErrored, atc.StatusFailed:
				errorJobs[job.ID] = job.FinishedBuild.ID
			default:
				break
			}
		}
	}
	return errorJobs
}

func filterRunningJobs(jobs []atc.Job) map[int]int {
	runningJobs := make(map[int]int)

	for _, job := range jobs {
		if job.NextBuild != nil && job.NextBuild.IsRunning() {
			runningJobs[job.ID] = job.NextBuild.ID
		}
	}

	return runningJobs
}
