package checker_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/SebastiaanPasterkamp/go-cache"
	discord "github.com/SebastiaanPasterkamp/go-discordwebhook"
	rss "github.com/mmcdole/gofeed"

	"github.com/SebastiaanPasterkamp/rsstodiscord/internal/checker"
)

func TestLoop(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testCases := []struct {
		name               string
		URL                string
		cached             string
		items              []*rss.Item
		expectedMessage    discord.Message
		expectedSendCalled bool
	}{
		{"Success", "http://example.com/feed", "other", []*rss.Item{&testItem},
			testMessage, true},
		{"Handle old results", "http://example.com/feed", testItem.GUID, []*rss.Item{&testItem},
			discord.Message{Content: "Unexpected"}, false},
		{"Handle no results", "http://example.com/feed", testItem.GUID, []*rss.Item{},
			discord.Message{Content: "Unexpected"}, false},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mem, err := cache.Configuration{
				InMemorySettings: &cache.InMemorySettings{},
			}.Factory(ctx)
			if err != nil {
				t.Fatalf("Failed to init cache: %v", err)
			}

			mem.Set(ctx, tt.cached, "sent", 5*time.Minute)

			var fetchCalled, sendCalled bool
			cfg := checker.Configuration{
				Fetch: func(feedURL string, ctx context.Context) (feed *rss.Feed, err error) {
					fetchCalled = true

					if feedURL != tt.URL {
						t.Errorf("Unexpected feed URL. Expected %q, got %q.",
							tt.URL, feedURL)
					}

					return &rss.Feed{
						Items: tt.items,
					}, nil
				},
				Send: func(msg discord.Message, wait bool) (*discord.Message, *discord.APIError, error) {
					sendCalled = true

					if !reflect.DeepEqual(msg, tt.expectedMessage) {
						t.Errorf("Unexpected message. Expected %v, got %v.",
							tt.expectedMessage, msg)

					}

					return &testMessage, nil, nil
				},
				Cache:    mem,
				URL:      tt.URL,
				Interval: 500 * time.Millisecond,
				Delay:    5 * time.Millisecond,
			}

			cctx, cancel := context.WithCancel(ctx)
			go func() {
				time.Sleep(1 * time.Second)
				defer cancel()
			}()

			cfg.Loop(cctx)

			if !fetchCalled {
				t.Error("Expected Fetch to be called but it was not.")
			}

			if !sendCalled && tt.expectedSendCalled {
				t.Error("Expected Send to be called but it was not.")
			}

			if sendCalled && !tt.expectedSendCalled {
				t.Error("Expected Send not to be called but it was.")
			}
		})
	}
}
