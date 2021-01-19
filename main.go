package main

import (
	"github.com/paroxity/portal/config"
	"github.com/paroxity/portal/query"
	"github.com/paroxity/portal/session"
	"github.com/paroxity/portal/socket"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"log"
	"strings"
)

func main() {
	l, err := minecraft.ListenConfig{
		AuthenticationDisabled: !config.Authentication(),
		StatusProvider:         query.StatusProvider{},
		ResourcePacks:          config.ResourcePacks(),
		TexturePacksRequired:   config.ForceTexturePacks(),
	}.Listen("raknet", ":19132")
	if err != nil {
		log.Fatalf("Unable to start listener: %v\n", err)
	}

	go func() {
		if err := socket.Listen(); err != nil {
			panic(err)
		}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Unable to accept connection: %v\n", err)
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
		log.Printf("%s failed to join: Server is whitelisted\n", conn.IdentityData().DisplayName)
		return
	}

	_, err := session.New(conn)
	if err != nil {
		log.Printf("Unable to create session, %v\n", err)
		_ = l.Disconnect(conn, text.Colourf("<red>%v</red>", err))
		return
	}
}
