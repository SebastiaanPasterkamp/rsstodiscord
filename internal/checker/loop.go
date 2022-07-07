package checker

import (
	"context"
	"errors"
	"log"
	"time"
)

// Loop continuously executes the Relay operation at the configured interval.
func (c Configuration) Loop(ctx context.Context) {
	log.Printf("Starting loop at %v interval with %v delay between items.",
		c.Interval, c.Delay)

	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()

	err := c.Relay(ctx)
	if err != nil && !errors.Is(err, ErrNoItems) {
		log.Printf("Failed to relay RSS to Discord: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := c.Relay(ctx)
			if err != nil && !errors.Is(err, ErrNoItems) {
				log.Printf("Failed to relay RSS to Discord: %v", err)
			}
		}
	}
}
