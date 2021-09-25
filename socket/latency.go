package socket

import (
	"github.com/paroxity/portal/internal"
	"github.com/paroxity/portal/session"
	"github.com/paroxity/portal/socket/packet"
	"time"
)

// ReportPlayerLatency sends the latency of each player to their connected server at the interval provided.
func (s *DefaultServer) ReportPlayerLatency(store session.Store, log internal.Logger, interval time.Duration) {
	for {
		for _, s := range store.All() {
			srv := s.Server()
			if srv == nil || !srv.Connected() {
				continue
			}
			if err := srv.Conn().WritePacket(&packet.UpdatePlayerLatency{
				PlayerUUID: s.UUID(),
				Latency:    s.Conn().Latency().Milliseconds(),
			}); err != nil {
				log.Errorf("failed to send packet: %v", err)
			}
		}
		time.Sleep(interval)
	}
}
