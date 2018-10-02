package concourse

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func TestReceiveRunningJobs(t *testing.T) {
	concourse := New(Config{
		Log: zerolog.New(os.Stdout),
		URL: "https://taa-ci-01.local.netconomy.net",
	})

	runningJobs, err := concourse.RunningJobs(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	t.Log(len(runningJobs))
}
