package checker_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/SebastiaanPasterkamp/go-cache"
	discord "github.com/SebastiaanPasterkamp/go-discordwebhook"
	rss "github.com/mmcdole/gofeed"

	"github.com/SebastiaanPasterkamp/rsstodiscord/internal/checker"
)

func TestRelay(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	testCases := []struct {
		name            string
		URL             string
		cached          string
		items           []*rss.Item
		fetchError      error
		sendError       error
		expectedMessage discord.Message
		expectedError   error
	}{
		{"Success", "http://example.com/feed", "other", []*rss.Item{&testItem}, nil, nil,
			testMessage, nil},
		{"Handle no results", "http://example.com/feed", "other", []*rss.Item{}, nil, nil,
			testMessage, checker.ErrNoItems},
		{"Handle old results", "http://example.com/feed", testItem.GUID, []*rss.Item{&testItem}, nil, nil,
			discord.Message{Content: "Unexpected"}, nil},
		{"Handle failed fetch", "whoops", "other", []*rss.Item{}, fmt.Errorf("bad url"), nil,
			discord.Message{Content: "Unexpected"}, checker.ErrFeedFailed},
		{"Handle failed send", "http://example.com/feed", "other", []*rss.Item{&testItem}, nil, fmt.Errorf("bad token"),
			testMessage, checker.ErrNotificationFailed},
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

			cfg := checker.Configuration{
				Fetch: func(feedURL string, ctx context.Context) (feed *rss.Feed, err error) {
					if feedURL != tt.URL {
						t.Errorf("Unexpected feed URL. Expected %q, got %q.",
							tt.URL, feedURL)
					}

					return &rss.Feed{
						Items: tt.items,
					}, tt.fetchError
				},
				Send: func(msg discord.Message, wait bool) (*discord.Message, *discord.APIError, error) {
					if !reflect.DeepEqual(msg, tt.expectedMessage) {
						t.Errorf("Unexpected message. Expected %v, got %v.",
							tt.expectedMessage, msg)

					}

					return &testMessage, nil, tt.sendError
				},
				Cache: mem,
				URL:   tt.URL,
				Delay: 5 * time.Millisecond,
			}

			err = cfg.Relay(ctx)
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Unexpected error. Expected %v, got %v.",
					tt.expectedError, err)
			}
		})
	}
}
