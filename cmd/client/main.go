package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
	"github.com/webtransport-example/config"
	"github.com/webtransport-example/logger"
)

const (
	ping = "Ping"
)

func main() {
	ctx := context.Background()
	log := logger.NewLog()
	cfg, err := config.New()
	if err != nil {
		log.Fatal("config", "New: %v", err)
	}

	client, err := NewDialer(cfg)
	if err != nil {
		log.Fatal("wt", "NewDialer: %v", err)
	}
	defer client.RoundTripper.Close()

	_, conn, err := client.Dial(ctx, cfg.Host, nil)
	if err != nil {
		log.Fatal("wt", "Dial: %v", err)
	}

	s, err := conn.OpenStream()
	if err != nil {
		log.Fatal("wt", "OpenStream: %v", err)
	}
	delay := time.Second * time.Duration(cfg.PingDelaySeconds)
	go pingSender(ctx, s, delay, log)

	for {
		buf := make([]byte, 1024)
		n, err := s.Read(buf)
		if err != nil {
			log.Error("wt", "Read: %v", err)
			break
		}
		incomingMessage := buf[:n]
		log.Debug("received", "Stream #%d Message: %s", s.StreamID(), incomingMessage)
	}
}

func NewDialer(cfg *config.Config) (*webtransport.Dialer, error) {
	var d webtransport.Dialer
	var qconf quic.Config
	var keyLog io.Writer
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	f, err := os.Create(cfg.WebTransportLogFile)
	if err != nil {
		return nil, err
	}
	keyLog = f
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: true,
			KeyLogWriter:       keyLog,
		},
		QuicConfig: &qconf,
	}
	d.RoundTripper = roundTripper
	return &d, nil
}

func pingSender(ctx context.Context, s webtransport.Stream, delay time.Duration, log logger.Logger) {
	timer := time.NewTimer(delay)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			outcomingMessage := ping
			_, err := s.Write([]byte(outcomingMessage))
			if err != nil {
				log.Error("wt", "Write: %v", err)
			}
			timer.Reset(delay)
			log.Debug("sent", "Stream #%d Message: %s", s.StreamID(), ping)
		}
	}
}
