package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	discord "github.com/SebastiaanPasterkamp/go-discordwebhook"
	"github.com/SebastiaanPasterkamp/gobernate"
	rss "github.com/mmcdole/gofeed"

	v "github.com/SebastiaanPasterkamp/rsstodiscord/internal/build"
	"github.com/SebastiaanPasterkamp/rsstodiscord/internal/checker"
	"github.com/SebastiaanPasterkamp/rsstodiscord/internal/cli"
)

func main() {
	ctx := context.Background()

	cfg, err := cli.Parse(os.Args)
	switch {
	case errors.Is(err, cli.ErrParsingFailed):
		log.Fatalf("Cannot parse CLI options: %v", err)
	case err != nil:
		os.Exit(1)
	}

	m, err := cfg.Cache.Factory(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}

	d, err := discord.New(cfg.Discord)
	if err != nil {
		log.Fatalf("Failed to initialize discord client: %v", err)
	}

	p := rss.NewParser()

	c := checker.Configuration{
		URL:      cfg.RSS,
		Cache:    m,
		TTL:      24 * time.Hour,
		Fetch:    p.ParseURLWithContext,
		Send:     d.Send,
		Timeout:  cfg.Timeout,
		Interval: cfg.Serve.Interval,
		Delay:    cfg.Delay,
	}

	switch {
	case cfg.Run != nil:
		err = c.Relay(ctx)
		if err != nil {
			log.Fatalf("Failed to relay RSS to Discord: %v", err)
		}
	case cfg.Serve != nil:
		g := gobernate.New(cfg.Serve.Port, v.Name, v.Version, v.Commit, v.Timestamp)

		go c.Loop(ctx)

		shutdown := g.Launch()
		defer g.Shutdown()

		g.Ready()
		<-shutdown
	default:
		log.Fatalf("Need to specify a subcommand. Use --help for information.")
	}

}
