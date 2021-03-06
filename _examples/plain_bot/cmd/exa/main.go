package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"

	"github.com/mavolin/adam/_examples/plain_bot/plugins/mod"
	"github.com/mavolin/adam/_examples/plain_bot/plugins/ping"
	"github.com/mavolin/adam/_examples/plain_bot/plugins/say"
	"github.com/mavolin/adam/pkg/bot"
	"github.com/mavolin/adam/pkg/impl/command/help"
)

func main() {
	b, err := bot.New(bot.Options{
		Token:            os.Getenv("DISCORD_BOT_TOKEN"),
		SettingsProvider: bot.NewStaticSettingsProvider(parsePrefixes()...),
		Owners:           parseOwners(),
		EditAge:          45 * time.Second,
		Status:           gateway.OnlineStatus,
		ActivityType:     discord.GameActivity,
		ActivityName:     "with the adam 🤖 framework",
	})
	if err != nil {
		log.Fatal(err)
	}

	addPlugins(b)

	log.Println("starting up")

	if err = b.Open(); err != nil {
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
	if len(prefixes) == 0 {
		return nil
	}

	return strings.Split(prefixes, ",")
}
