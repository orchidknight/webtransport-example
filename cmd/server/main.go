package main

import (
	"context"
	"net/http"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
	"github.com/webtransport-example/config"
	"github.com/webtransport-example/logger"
)

const (
	ping = "Ping"
	pong = "Pong"
)

func main() {
	ctx := context.Background()
	log := logger.NewLog()
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config", "New: %v", err)
	}

	s := webtransport.Server{
		H3: http3.Server{Addr: cfg.Port},
	}

	http.HandleFunc("/wt", func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.Upgrade(w, r)
		if err != nil {
			log.Error("wt", "Upgrade: upgrading failed: %s", err)
			w.WriteHeader(500)
			return
		}

		s, err := conn.AcceptStream(ctx)
		log.Debug("wt", "Stream #%d was accepted.", s.StreamID())
		if err != nil {
			log.Error("wt", "AcceptStream: %v", err)
		}
		for {
			buf := make([]byte, 1024)
			n, err := s.Read(buf)
			if err != nil {
				break
			}
			incomingMessage := string(buf[:n])
			outcomingMessage := pong
			log.Debug("received", "Stream #%d Message: %s", s.StreamID(), incomingMessage)
			if incomingMessage == ping {
				_, err = s.Write([]byte(outcomingMessage))
				if err != nil {
					log.Error("wt", "Write: %v", err)

				}
				log.Debug("sent", "Stream #%d Message: %s", s.StreamID(), outcomingMessage)
			}
		}

	})

	err = s.ListenAndServeTLS(cfg.CertificatePath, cfg.KeyPath)
	if err != nil {
		log.Fatal("server", "ListenAndServeTLS: %v", err)
	}
}
