package socket

import (
	"github.com/paroxity/portal/session"
	"github.com/paroxity/portal/socket/packet"
	"github.com/sirupsen/logrus"
	"time"
)

// ReportPlayerLatency sends the latency of each player to their connected server at the interval provided.
// TODO: Redesign latency reporting to become more modular.
func ReportPlayerLatency(interval time.Duration) {
	for {
		for _, s := range session.All() {
			srv := s.Server()
			if srv == nil || !srv.Connected() {
				continue
			}
			if err := srv.Conn().WritePacket(&packet.UpdatePlayerLatency{
				PlayerUUID: s.UUID(),
				Latency:    s.Conn().Latency().Milliseconds(),
			}); err != nil {
				logrus.Error(err)
			}
		}
		time.Sleep(interval)
	}
}
