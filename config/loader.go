package config

import (
	"fmt"
	"github.com/paroxity/wormhole/server"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type config struct {
	Query struct {
		MaxPlayers int
		MOTD       string
	}
	Whitelist struct {
		Enabled bool
		Players []string
	}
	Proxy struct {
		Groups       []string
		DefaultGroup string
		Servers      []ServerInfo
	}
}

type ServerInfo struct {
	Name    string
	Group   string
	Address string
}

func init() {
	if err := Load(); err != nil {
		log.Printf("Unable to load config: %v\n", err)
	}
}

func Load() error {
	c := config{}
	c.Query.MOTD = "Wormhole"
	c.Proxy.Groups = []string{"Hub"}
	c.Proxy.DefaultGroup = "Hub"
	c.Proxy.Servers = append(c.Proxy.Servers, ServerInfo{
		Name:    "Hub1",
		Group:   "Hub",
		Address: "127.0.0.1",
	})

	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		data, err := yaml.Marshal(c)
		if err != nil {
			return fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := ioutil.WriteFile("config.yml", data, 0644); err != nil {
			return fmt.Errorf("failed creating config: %v", err)
		}
	} else {
		data, err := ioutil.ReadFile("config.yml")
		if err != nil {
			return fmt.Errorf("error reading config: %v", err)
		}
		if err := yaml.Unmarshal(data, &c); err != nil {
			return fmt.Errorf("error decoding config: %v", err)
		}
	}

	maxPlayers = c.Query.MaxPlayers
	motd = c.Query.MOTD
	whitelist = c.Whitelist.Enabled
	whitelisted = c.Whitelist.Players

	for _, name := range c.Proxy.Groups {
		server.AddGroup(server.NewGroup(name))
	}

	g, ok := server.GroupFromName(c.Proxy.DefaultGroup)
	if !ok {
		panic(fmt.Sprintf("default group %s not found", c.Proxy.DefaultGroup))
	}
	server.SetDefaultGroup(g)

	for _, info := range c.Proxy.Servers {
		if _, err := server.New(info.Name, info.Group, info.Address); err != nil {
			panic(err)
		}
	}

	return nil
}
