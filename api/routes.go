package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/chi"
)

// Server represents all state required to serve api requests.
// Server ought to live for the entirety of the application's expected run time
type Server struct {
	// Note that the below session may not belong exclusively to server in the future
	// Discord sessions should be treated as "write only" by the time they reach server.
	Discord      *discordgo.Session
	AlertChannel *discordgo.Channel
}

func (s *Server) ConfigureRoutes(r chi.Router) {
	r.Post("/alert", s.alert)
}
