package main

import (
	"github.com/paroxity/portal"
	"github.com/paroxity/portal/log"
	"github.com/paroxity/portal/server"
	"github.com/paroxity/portal/session"
	"github.com/paroxity/portal/socket"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
)

func main() {
	c := portal.DefaultConfig() // TODO: Load from file

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
	if c.Logger.File != "" {
		fileLogger, err := log.New(c.Logger.File)
		if err != nil {
			logger.Fatalf("unable to create file logger: %v", err)
		}
		logger.SetOutput(fileLogger)
	}
	level, err := logrus.ParseLevel(c.Logger.Level)
	if err != nil {
		logger.Errorf("unable to parse log level '%s': %v", c.Logger.Level, err)
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	serverRegistry := server.NewDefaultRegistry()
	p := portal.Options{
		Logger: logger,

		Address: c.Network.Address,
		ListenConfig: minecraft.ListenConfig{
			StatusProvider: portal.NewMOTDStatusProvider("Portal"),
		},

		ServerRegistry: serverRegistry,
		LoadBalancer:   session.NewSplitLoadBalancer(serverRegistry),
	}.Portal()
	if err := p.Listen(); err != nil {
		logger.Fatalf("failed to listen on %s: %v", c.Network.Address, err)
	}

	socketServer := socket.NewDefaultServer(c.Network.Communication.Address, c.Network.Communication.Secret, serverRegistry, logger)
	if err := socketServer.Listen(); err != nil {
		p.Logger().Fatalf("socket server failed to listen: %v", err)
	}

	for {
		s, err := p.Accept()
		if err != nil {
			p.Logger().Errorf("failed to accept connection for %s: %v", s.Conn().IdentityData().DisplayName, err)
			_ = p.Disconnect(s.Conn(), text.Colourf("<red>%v</red>", err))
			continue
		}
		p.Logger().Infof("%s has been connected to server %s", s.Conn().IdentityData().DisplayName, s.Server().Name())
	}
}
