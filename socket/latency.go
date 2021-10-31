package socket

import (
	"github.com/paroxity/portal/socket/packet"
	"time"
)

// ReportPlayerLatency sends the latency of each player to their connected server at the interval provided.
func (s *DefaultServer) ReportPlayerLatency(interval time.Duration) {
	for {
		for _, session := range s.SessionStore().All() {
			srv := session.Server()
			if srv == nil || !srv.Connected() {
				continue
			}
			if err := srv.Conn().WritePacket(&packet.UpdatePlayerLatency{
				PlayerUUID: session.UUID(),
				Latency:    session.Conn().Latency().Milliseconds(),
			}); err != nil {
				s.Logger().Errorf("failed to send packet: %v", err)
			}
		}
		time.Sleep(interval)
	}
}
