package checker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SebastiaanPasterkamp/go-cache"
)

// Relay fetches the configured RSS feed and any new items
// to the Discord Webhook client.
func (c Configuration) Relay(ctx context.Context) error {
	log.Printf("Fetching %q.\n", c.URL)

	feed, err := c.Fetch(c.URL, ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFeedFailed, err)
	}

	if len(feed.Items) == 0 {
		return ErrNoItems
	}

	ticker := time.NewTicker(c.Delay)
	defer ticker.Stop()

	for i := len(feed.Items) - 1; i >= 0; i-- {
		item := feed.Items[i]

		select {
		case <-ctx.Done():
			return ErrContextCancelled
		case <-ticker.C:
			mem := ""
			err := c.Cache.Get(ctx, item.GUID, &mem)
			switch err {
			case nil:
				// Already seen
				continue
			case cache.ErrNotFound:
				// New item, resume function
			default:
				return fmt.Errorf("%w, %v", ErrCacheFailed, err)
			}

			m := Translate(item)

			_, fail, err := c.Send(m, true)
			if err != nil {
				return fmt.Errorf("%w: %v: %v", ErrNotificationFailed, err, fail)
			}

			if fail != nil {
				return fmt.Errorf("%w: %v: %v", ErrNotificationFailed, err, fail)
			}

			err = c.Cache.Set(ctx, item.GUID, "sent", c.TTL)
			if err != nil {
				return fmt.Errorf("%w, %v", ErrCacheFailed, err)
			}

			log.Println(m.Embeds[0].Title)
		}
	}

	return nil
}
