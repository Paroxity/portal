package main

import (
	"encoding/json"
	"github.com/paroxity/portal"
	portallog "github.com/paroxity/portal/log"
	"github.com/paroxity/portal/session"
	"github.com/paroxity/portal/socket"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	conf := readConfig()

	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})
	if conf.Logger.File != "" {
		fileLogger, err := portallog.New(conf.Logger.File)
		if err != nil {
			logger.Fatalf("unable to create file logger: %v", err)
		}
		logger.SetOutput(fileLogger)
	}
	level, err := logrus.ParseLevel(conf.Logger.Level)
	if err != nil {
		logger.Errorf("unable to parse log level '%s': %v", conf.Logger.Level, err)
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	resourcePacks, err := portal.LoadResourcePacks(conf.ResourcePacks.Directory)
	if err != nil {
		logger.Fatalf("unable to load resource packs: %v", err)
	}

	p := portal.New(portal.Options{
		Logger: logger,

		Address: conf.Network.Address,
		ListenConfig: minecraft.ListenConfig{
			StatusProvider: portal.NewMOTDStatusProvider("Portal"),

			ResourcePacks:        resourcePacks,
			TexturePacksRequired: conf.ResourcePacks.Required,
		},

		Whitelist: session.NewSimpleWhitelist(conf.Whitelist.Enabled, conf.Whitelist.Players),
	})
	if err := p.Listen(); err != nil {
		logger.Fatalf("failed to listen on %s: %v", conf.Network.Address, err)
	}

	socketServer := socket.NewDefaultServer(conf.Network.Communication.Address, conf.Network.Communication.Secret, p.SessionStore(), p.ServerRegistry(), logger)
	if err := socketServer.Listen(); err != nil {
		p.Logger().Fatalf("socket server failed to listen: %v", err)
	}
	if conf.PlayerLatency.Report {
		go socketServer.ReportPlayerLatency(time.Second * time.Duration(conf.PlayerLatency.UpdateInterval))
	}

	for {
		s, err := p.Accept()
		if err != nil {
			p.Logger().Errorf("failed to accept connection: %v", err)
			continue
		}
		p.Logger().Infof("%s has been connected to server %s", s.Conn().IdentityData().DisplayName, s.Server().Name())
	}
}

func readConfig() portal.Config {
	c := portal.DefaultConfig()
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		f, err := os.Create("config.json")
		if err != nil {
			log.Fatalf("error creating config: %v", err)
		}
		data, err := json.MarshalIndent(c, "", "\t")
		if err != nil {
			log.Fatalf("error encoding default config: %v", err)
		}
		if _, err := f.Write(data); err != nil {
			log.Fatalf("error writing encoded default config: %v", err)
		}
		_ = f.Close()
	}
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	if err := json.Unmarshal(data, &c); err != nil {
		log.Fatalf("error decoding config: %v", err)
	}
	return c
}
