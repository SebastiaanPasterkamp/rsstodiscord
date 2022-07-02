package checker

// This package translates mmcdole/gofeed.Feed model into a discord.Message.

import (
	discord "github.com/SebastiaanPasterkamp/go-discordwebhook"
	rss "github.com/mmcdole/gofeed"
)

// Translate translates a mcdole.gofeed.Item into a discord.Message
func Translate(i *rss.Item) discord.Message {
	return discord.Message{
		Content:   "New Raspberry Pi in stock",
		Username:  "rpilocator",
		AvatarURL: "https://rpilocator.com/favicon.png",
		Embeds: []discord.Embed{
			{
				Title:       i.Title,
				URL:         i.Link,
				Type:        discord.EmbedRich,
				Description: i.Description,
				Timestamp:   *i.PublishedParsed,
				Color:       discord.ColorBlue,
			},
		},
		AllowedMentions: discord.AllowedMention{
			Parse: []discord.MentionType{},
		},
	}
}
