package checker

import (
	"context"
	"time"

	"github.com/SebastiaanPasterkamp/go-cache"
	discord "github.com/SebastiaanPasterkamp/go-discordwebhook"
	rss "github.com/mmcdole/gofeed"
)

// Configuration offers settings to start an rss feed
// check routine.
type Configuration struct {
	// URL is the RSS feed location to check.
	URL string
	// Cache is a means to remember which items have already been relayed.
	Cache cache.Repository
	// TTL is the duration for the cache to persist. Items still in the feed
	// after this time will be resent as "new" items.
	TTL time.Duration
	// Fetch takes an RSS URL and returns feed an rss.Feed or an error.
	Fetch func(feedURL string, ctx context.Context) (feed *rss.Feed, err error)
	// Send takes a discordwebhook.Message and sends this to discord.
	Send func(m discord.Message, wait bool) (*discord.Message, *discord.APIError, error)
	// Timeout limits how long the rss feed poll can take.
	Timeout time.Duration
	// Interval determines how often the feed is going to
	// be polled.
	Interval time.Duration
	// Delay reduces the rate at which new items are relayed to Discord. This
	// prevents burst spamming a channel and avoids 429 responses.
	Delay time.Duration
}
