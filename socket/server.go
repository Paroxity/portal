package socket

import (
	"github.com/paroxity/portal/config"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
)

// Listen starts a TCP listener on the configured address to handle incoming connections.
func Listen() error {
	listener, err := net.Listen("tcp", config.SocketAddress())
	if err != nil {
		return err
	}
	logrus.Infof("Socket server listening on %s\n", config.SocketAddress())

	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Infoln(err)
			continue
		}
		logrus.Debugln("Socket server accepted a new connection")

		go handleClient(NewClient(conn))
	}
}

// handleClient reads packets from a connected client and handles them with the appropriate handler.
func handleClient(c *Client) {
	defer c.Close()

	for {
		pk, err := c.ReadPacket()
		if err != nil {
			if containsAny(err.Error(), "EOF", "closed") {
				return
			}
			logrus.Errorln(err)
			continue
		}

		h, ok := handlers[pk.ID()]
		if ok {
			if err := h.Handle(pk, c); err != nil {
				logrus.Errorln(err)
			}
		} else {
			if c.name == "" {
				logrus.Debugf("Unhandled packet %T from unauthorized socket connection\n", pk)
			} else {
				logrus.Debugf("Unhandled packet %T from %s socket connection\n", pk, c.name)
			}
		}
	}
}

// containsAny checks if the string contains any of the provided sub strings.
func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}

	return false
}
