package cli

import (
	"fmt"
	"time"

	"github.com/SebastiaanPasterkamp/go-cache"

	"github.com/SebastiaanPasterkamp/rsstodiscord/internal/build"
)

// Configuration is the data structure for the CLI.
type Configuration struct {
	// Cache ensures only new items are relayed
	Cache cache.Configuration
	// Base is the part of the configuration that args can handle
	Base
}

// Base is the part of the configuration that args can handle
type Base struct {
	// Run is a collection of settings to perform a one-off execution with.
	Run *Run `arg:"subcommand:run" help:"One-off execution of the RSS discord relay."`
	// Serve is a collection of settings to launch a Kubernetes service which
	// relays messages at a specified interval.
	Serve *Serve `arg:"subcommand:serve" help:"Launch a server to continuously relay RSS to discord."`
	// Discord is the discord webhook URL to relay messages to.
	Discord string `arg:"--discord-webhook,env:DISCORD_WEBHOOK,required" placeholder:"https://discord.com/api/webhooks/.../..." help:"Secret URL to the Discord Webhook API."`
	// RSS is the URL to monitor
	RSS string `arg:"--rss-url,env:RSS_URL,required" placeholder:"http://example.com/feed" help:"URL of the RSS feed to fetch."`
	// Timeout limits how long the rss feed poll can take.
	Timeout time.Duration `args:"--timeout" default:"30s" help:"Timeout on fetching the RSS feed."`
	// Delay reduces the rate at which new items are relayed to Discord. This
	// prevents burst spamming a channel and avoids 429 responses.
	Delay time.Duration `args:"--delay" default:"5s" help:"Delay between relaying newly found items."`
}

// Run describes a CLI subcommand to do a one-time execution
type Run struct {
}

// Serve describes a CLI subcommand to launch a service
type Serve struct {
	Port string `args:"--port,env:PORT" default:"8080"`
	// Interval determines how often the feed is going to
	// be polled.
	Interval time.Duration `args:"--interval" default:"5m" help:"RSS polling interval."`
}

// Arguments is the data structure for the CLI
type Arguments struct {
	Base
	cache.InMemorySettings
	cache.RedisSettings
}

// Version provides the CLI help with the current version of the binary.
func (Arguments) Version() string {
	return fmt.Sprintf("%s %s (%s-%s @ %s)\n",
		build.Name, build.Version, build.Branch, build.Commit, build.Timestamp)
}

// Description provides the CLI help with a description of the binary.
func (Arguments) Description() string {
	return `This program parses an RSS feed, converts them into messages, and posts them to
a Discord channel using a Discord webhook. The RSS item IDs are stored in cache,
so future iterations only pick up new items.`
}
