package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	"github.com/haashi/omega-strikers-bot/internal/match"
	"github.com/haashi/omega-strikers-bot/internal/matchmaking"
	"github.com/haashi/omega-strikers-bot/internal/slashcommands"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() { flag.Parse() }

func main() {
	rand.Seed(time.Now().Unix())
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	logLevel := os.Getenv("loglevel")
	if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	db.Init()
	match.Init()
	matchmaking.Init()
	discord.Init()
	slashcommands.Init()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("press Ctrl+C to exit")
	<-stop
	slashcommands.Stop()
	discord.Stop()

	log.Println("gracefully shutting down.")
}