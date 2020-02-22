package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Alert struct {
	Source   string `json:"source"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

func (s *Server) alert(w http.ResponseWriter, r *http.Request) {
	var alert Alert

	buffer, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Printf("Could not read request body due to error %s\n", err.Error())
		w.WriteHeader(500)
		return
	}

	err = json.Unmarshal(buffer, &alert)

	if err != nil {
		log.Printf("Could not unmarshal request JSON, message will not be submitted %s\n", err.Error())
		log.Printf("Submitted json was %s\n", string(buffer))
		w.WriteHeader(400)
		return
	}

	go func() {
		_, err = s.Discord.ChannelMessageSend(s.AlertChannel.ID, fmt.Sprintf("Alert from: %s\nSeverity: %s\nMessage: %s", alert.Source, alert.Severity, alert.Message))

		if err != nil {
			log.Printf("Failed to send channel message due to error %s\n", err.Error())
		}
	}()

	w.WriteHeader(200)
}
