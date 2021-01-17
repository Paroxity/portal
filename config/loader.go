package config

import (
	"fmt"
	"github.com/paroxity/portal/server"
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path"
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
		Groups         map[string]map[string]string
		DefaultGroup   string
		Authentication bool
		ResourcesDir   string
		ForceTextures  bool
	}
	Socket struct {
		BindAddress string
		Secret      string
	}
}

func init() {
	if err := Load(); err != nil {
		log.Printf("Unable to load config: %v\n", err)
	}
}

func Load() error {
	c := config{}
	c.Query.MOTD = "Portal"
	c.Proxy.Groups = map[string]map[string]string{
		"Hub": {
			"Hub1": "127.0.0.1:19133",
		},
	}
	c.Proxy.DefaultGroup = "Hub"
	c.Proxy.Authentication = true
	c.Socket.BindAddress = "127.0.0.1:19131"

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

	if len(c.Proxy.Groups) == 0 {
		return fmt.Errorf("groups are empty")
	}

	for name, group := range c.Proxy.Groups {
		if len(group) == 0 {
			return fmt.Errorf("group %s has no servers", group)
		}

		g := server.NewGroup(name)
		for name, address := range group {
			g.AddServer(server.New(name, address))
		}

		server.AddGroup(g)
	}

	g, ok := server.GroupFromName(c.Proxy.DefaultGroup)
	if !ok {
		return fmt.Errorf("default group %s not found", c.Proxy.DefaultGroup)
	}
	server.SetDefaultGroup(g)

	authentication = c.Proxy.Authentication
	files, err := ioutil.ReadDir(c.Proxy.ResourcesDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, file := range files {
		pack, err := resource.Compile(path.Join(c.Proxy.ResourcesDir, file.Name()))
		if err != nil {
			return err
		}
		resourcePacks = append(resourcePacks, pack)
	}
	forceTextures = c.Proxy.ForceTextures

	socketAddress = c.Socket.BindAddress
	socketSecret = c.Socket.Secret

	return nil
}
