package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/mavolin/adam/examples/plain_bot/plugins/embed"
	"github.com/mavolin/adam/examples/plain_bot/plugins/mod"
	"github.com/mavolin/adam/examples/plain_bot/plugins/ping"
	"github.com/mavolin/adam/examples/plain_bot/plugins/say"
	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/impl/command/help"
)

func main() {
	b, err := bot.New(bot.Options{
		Token:            os.Getenv("DISCORD_BOT_TOKEN"),
		SettingsProvider: bot.StaticSettings(parsePrefixes()...),
		Owners:           parseOwners(),
		EditAge:          45 * time.Second,
		ActivityName:     "with the adam 🤖 framework",
	})
	if err != nil {
		log.Fatal(err)
	}

	addPlugins(b)

	log.Println("starting up")

	if err = b.Open(2 * time.Second); err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig

	log.Println("received SIGINT, shutting down")

	if err = b.Close(); err != nil {
		log.Println("could not close bot properly:", err.Error())
	}
}

func addPlugins(b *bot.Bot) {
	b.AddCommand(help.New(help.Options{}))

	b.AddCommand(embed.New())
	b.AddCommand(ping.New())
	b.AddCommand(say.New())

	b.AddModule(mod.New())
}

func parseOwners() []discord.UserID {
	rawOwners := strings.Split(os.Getenv("BOT_OWNERS"), ",")

	owners := make([]discord.UserID, 0, len(rawOwners))

	for _, o := range rawOwners {
		s, err := discord.ParseSnowflake(o)
		if err == nil {
			owners = append(owners, discord.UserID(s))
		}
	}

	return owners
}

func parsePrefixes() []string {
	prefixes := os.Getenv("BOT_PREFIXES")
	if prefixes == "" {
		return nil
	}

	return strings.Split(prefixes, ",")
}
