package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/chaos-ops/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

func main() {
	// Configuring logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Handling environment variables
	viper.SetDefault("RequestTimeout", 10)
	viper.SetDefault("ListenAddress", ":8080")
	viper.SetDefault("BehindProxy", false)

	viper.BindEnv("DiscordAPIKey", "CHAOS_DISCORD_API_KEY")
	viper.BindEnv("DiscordAlertChannel", "CHAOS_DISCORD_ALERT_CHANNEL")
	viper.BindEnv("RequestTimeout", "CHAOS_REQUEST_TIMEOUT")
	viper.BindEnv("ListenAddress", "CHAOS_LISTEN_ADDRESS")
	viper.BindEnv("BehindProxy", "CHAOS_BEHIND_PROXY")

	key := viper.GetString("DiscordAPIKey")
	channel := viper.GetString("DiscordAlertChannel")
	timeout := viper.GetDuration("RequestTimeout")
	address := viper.GetString("ListenAddress")
	behindProxy := viper.GetBool("BehindProxy")

	if key == "" {
		log.Fatalf("Cannot start server with empty discord api key")
	}
	if channel == "" {
		log.Fatalf("Cannot start server with empty alert channel")
	}

	// Setting up discordGo
	dg, err := discordgo.New("Bot " + key)

	if err != nil {
		log.Fatalf("Could not create discord instance due to error %s", err.Error())
	}

	st, err := dg.Channel(channel)

	if err != nil {
		log.Fatalf("Could not resolve channel due to error %s", err.Error())
	}

	server := api.Server{
		Discord:      dg,
		AlertChannel: st,
	}

	// Setting up chi
	r := chi.NewRouter()
	if behindProxy {
		r.Use(middleware.RealIP)
	}
	r.Use(middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(timeout*time.Second),
		middleware.Heartbeat("/health"))
	r.Group(server.ConfigureRoutes)

	log.Printf("Starting http server at address %s...", address)

	err = http.ListenAndServe(address, r)

	if err != nil {
		log.Fatalf("Cannot start server due to error %s", err.Error())
	}
}
