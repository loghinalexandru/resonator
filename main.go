package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/internal/command"
	"github.com/loghinalexandru/resonator/pkg/logging"
)

var (
	token        string
	swearsApiURL string
	logLevel     logging.LogLevel = logging.Info
	shardID      int              = 0
	shardCount   int              = 1
)

func loadEnv() {
	token = os.Getenv("BOT_TOKEN")
	swearsApiURL = os.Getenv("SWEARS_API_URL")

	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		logLevel = logging.ToLogLevel(lvl)
	}

	if replicas := os.Getenv("SHARD_COUNT"); replicas != "" {
		shardCount, _ = strconv.Atoi(replicas)
	}

	if replicaID := os.Getenv("SHARD_ID"); replicaID != "" {
		index := strings.Split(replicaID, "-")
		shardID, _ = strconv.Atoi(index[len(index)-1])
	}
}

func getIntents() discordgo.Intent {
	return discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
}

func main() {
	loadEnv()

	session, sessionError := discordgo.New("Bot " + token)
	session.ShouldReconnectVoiceConnOnError = false

	logger := logging.New(logLevel, log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile))

	cmdSync := sync.Map{}
	cmds := []CustomCommandDef{
		command.NewPlay(&cmdSync),
		command.NewReact(&cmdSync),
		command.NewRo(&cmdSync),
		command.NewCurse(&cmdSync, swearsApiURL),
		command.NewFeed(&cmdSync, swearsApiURL),
		command.NewSwear(swearsApiURL),
		command.NewAnime(),
		command.NewManga(),
		command.NewQuote(),
	}

	handlers := []any{
		Join(logger),
		InteractionCreate(cmds, logger),
	}

	if sessionError != nil {
		logger.Error(sessionError)
		return
	}

	for _, handler := range handlers {
		session.AddHandler(handler)
	}

	session.Identify.Intents = getIntents()
	session.ShardID = shardID
	session.ShardCount = shardCount

	socketError := session.Open()
	defer session.Close()

	if socketError != nil {
		logger.Error(socketError)
		return
	}

	logger.Info("Bot is ready!")
	logger.Info("Bot ShardId: ", session.ShardID)
	logger.Info("Bot ShardCount: ", session.ShardCount)

	for _, command := range cmds {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, "", command.Definition())

		if err != nil {
			logger.Error(err)
		}
	}

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigTerm
}
