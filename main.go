package main

import (
	"flag"
	"fmt"
	"github.com/paroxity/portal/config"
	"github.com/paroxity/portal/logger"
	"github.com/paroxity/portal/query"
	"github.com/paroxity/portal/session"
	"github.com/paroxity/portal/socket"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"strings"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.json", "Path to the config file.")
	flag.Parse()

	if err := config.Load(configPath); err != nil {
		fmt.Printf("Unable to load config: %v", err)
	}
	log, err := logger.New(config.LogFile())
	if err != nil {
		panic(err)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
	logrus.SetOutput(log)
	logrus.SetLevel(config.LogLevel())

	l, err := minecraft.ListenConfig{
		AuthenticationDisabled: !config.Authentication(),
		StatusProvider:         query.StatusProvider{},
		ResourcePacks:          config.ResourcePacks(),
		TexturePacksRequired:   config.ForceTexturePacks(),
	}.Listen("raknet", config.BindAddress())
	if err != nil {
		logrus.Fatalf("Unable to start listener: %v\n", err)
	}
	logrus.Infof("Listening on %s\n", config.BindAddress())

	go func() {
		if err := socket.Listen(); err != nil {
			panic(err)
		}
	}()
	if config.ReportPlayerLatency() {
		go socket.ReportPlayerLatency(config.PlayerLatencyUpdateInterval())
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			logrus.Infof("Unable to accept connection: %v\n", err)
			return
		}

		go handleConnection(l, conn.(*minecraft.Conn))
	}
}

// handleConnection handles an incoming connection from the Listener.
func handleConnection(l *minecraft.Listener, conn *minecraft.Conn) {
	var whitelisted bool
	for _, p := range config.Whitelisted() {
		if strings.EqualFold(conn.IdentityData().DisplayName, p) {
			whitelisted = true
			break
		}
	}
	if config.Whitelist() && !whitelisted {
		_ = l.Disconnect(conn, text.Colourf("<red>Server is whitelisted</red>"))
		logrus.Infof("%s failed to join: Server is whitelisted\n", conn.IdentityData().DisplayName)
		return
	}

	s, err := session.New(conn)
	if err != nil {
		logrus.Errorf("Unable to create session, %v\n", err)
		_ = l.Disconnect(conn, text.Colourf("<red>%v</red>", err))
		return
	}
	logrus.Infof("%s has been connected to server %s in group %s\n", s.Conn().IdentityData().DisplayName, s.Server().Name(), s.Server().Group())
}
