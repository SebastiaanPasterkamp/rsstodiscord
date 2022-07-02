package checker

import "fmt"

var (
	// ErrFeedFailed is the error returned when the FeedParser did not return
	// successfully.
	ErrFeedFailed = fmt.Errorf("fetching rss failed")
	// ErrNoItems is the error returned when the FeedParser was successful, did
	// not yield any items.
	ErrNoItems = fmt.Errorf("no items")
	// ErrCacheFailed is the error returned when the cache could either not
	// retrieve or store the item ID.
	ErrCacheFailed = fmt.Errorf("cache failed")
	// ErrNotificationFailed is the error returned when the new rss feed item
	// could not be sent to Discord.
	ErrNotificationFailed = fmt.Errorf("discord notification failed")
	// ErrContextCancelled is the error returned when the context is cancelled
	// before all items have been processed.
	ErrContextCancelled = fmt.Errorf("context cancelled")
)
