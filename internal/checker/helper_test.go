package checker_test

import (
	"time"

	discord "github.com/SebastiaanPasterkamp/go-discordwebhook"
	rss "github.com/mmcdole/gofeed"
)

var testDate, _ = time.Parse("Mon 02 Jan 2006 15:04:05 GMT", "Sat 02 Jul 2022 14:23:46 GMT")

var testItem = rss.Item{
	Title:           "Stock Alert (NL): RPi 4 Model B - 2GB RAM is In Stock at Elektor",
	Description:     "Stock Alert (NL): RPi 4 Model B - 2GB RAM is In Stock at Elektor",
	Link:            "https://rpilocator.com?vendor=elektor&utm_source=feed&utm_medium=rss",
	Links:           []string{"https://rpilocator.com?vendor=elektor&utm_source=feed&utm_medium=rss"},
	Published:       "Sat 02 Jul 2022 14:23:46 GMT",
	PublishedParsed: &testDate,
	GUID:            "6AB88F05-64BB-4BDC-A27B816617C971A1",
	Categories:      []string{"elektor", "NL", "PI4"},
}

var testMessage = discord.Message{
	Content:   "New Raspberry Pi in stock",
	Username:  "rpilocator",
	AvatarURL: "https://rpilocator.com/favicon.png",
	Embeds: []discord.Embed{
		{
			Title:       "Stock Alert (NL): RPi 4 Model B - 2GB RAM is In Stock at Elektor",
			URL:         "https://rpilocator.com?vendor=elektor&utm_source=feed&utm_medium=rss",
			Type:        discord.EmbedRich,
			Description: "Stock Alert (NL): RPi 4 Model B - 2GB RAM is In Stock at Elektor",
			Timestamp:   testDate,
			Color:       discord.ColorBlue,
		},
	},
	AllowedMentions: discord.AllowedMention{
		Parse: []discord.MentionType{},
	},
}
