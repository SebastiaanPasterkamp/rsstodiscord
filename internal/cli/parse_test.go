package cli_test

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/SebastiaanPasterkamp/go-cache"
	"github.com/SebastiaanPasterkamp/rsstodiscord/internal/cli"
	"github.com/alexflint/go-arg"
)

func TestColorToString(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expectedError  error
		expectedConfig *cli.Configuration
	}{
		{"Minimum args gives in-memory", []string{"rsstodiscord", "run", "--discord-webhook", "foo", "--rss-url", "bar"},
			nil, &cli.Configuration{
				Cache: cache.Configuration{
					InMemorySettings: &cache.InMemorySettings{},
				},
				Base: cli.Base{
					Discord: "foo",
					RSS:     "bar",
					Run:     &cli.Run{},
					Timeout: 30 * time.Second,
					Delay:   5 * time.Second,
				},
			}},
		{"Adding cache uses redis", []string{"rsstodiscord", "run", "--discord-webhook", "foo", "--rss-url", "bar", "--redis-address", "ruh"},
			nil, &cli.Configuration{
				Cache: cache.Configuration{
					RedisSettings: &cache.RedisSettings{
						Address: "ruh",
					},
				},
				Base: cli.Base{
					Discord: "foo",
					RSS:     "bar",
					Run:     &cli.Run{},
					Timeout: 30 * time.Second,
					Delay:   5 * time.Second,
				},
			}},
		{"Show help", []string{"rsstodiscord", "--help"},
			arg.ErrHelp, nil},
		{"Show version", []string{"rsstodiscord", "--version"},
			arg.ErrVersion, nil},
		{"Require subcommand", []string{"rsstodiscord", "--discord-webhook", "foo", "--rss-url", "bar"},
			arg.ErrHelp, nil},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			backup := os.Args
			defer func() {
				os.Args = backup
			}()
			os.Args = tt.args

			result, err := cli.Parse(tt.args)
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Unexpected error. Expected %v, got %v",
					tt.expectedError, err)
			}

			exp, _ := json.Marshal(tt.expectedConfig)
			res, _ := json.Marshal(result)

			if string(exp) != string(res) {
				t.Errorf("Unexpected configuration. Expected %v, got %v.",
					string(exp), string(res))
			}
		})
	}
}
