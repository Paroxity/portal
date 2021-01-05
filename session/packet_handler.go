package session

import (
	"fmt"
	"github.com/paroxity/portal/server"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"log"
	"strings"
)

// handlePackets handles the packets sent between the client and the server. Processes such as runtime
// translations are also handled here.
func handlePackets(s *Session) {
	go func() {
		defer func() {
			s.Close()
		}()
		for {
			pk, err := s.Conn().ReadPacket()
			if err != nil {
				log.Println(err)
				return
			}
			s.translatePacket(pk)

			switch pk := pk.(type) {
			case *packet.CommandRequest:
				args := strings.Split(pk.CommandLine, " ")
				switch args[0][1:] {
				case "server":
					_ = s.conn.WritePacket(&packet.Text{
						TextType: packet.TextTypeRaw,
						Message:  text.Colourf("<green>You are currently on %s</green>", s.server.Name()),
					})
					continue
				case "transfer":
					if len(args) < 3 {
						_ = s.conn.WritePacket(&packet.Text{
							TextType: packet.TextTypeRaw,
							Message:  text.Colourf("<red>Please provide a group and a server to transfer to</red>"),
						})
						continue
					}
					g, ok := server.GroupFromName(args[1])
					if !ok {
						_ = s.conn.WritePacket(&packet.Text{
							TextType: packet.TextTypeRaw,
							Message:  text.Colourf("<red>Group %s not found</red>", args[1]),
						})
						continue
					}
					srv, ok := g.Server(args[2])
					if !ok {
						_ = s.conn.WritePacket(&packet.Text{
							TextType: packet.TextTypeRaw,
							Message:  text.Colourf("<red>Server %s not found in group %s</red>", args[2], g.Name()),
						})
						continue
					}
					if err := s.Transfer(srv); err != nil {
						_ = s.conn.WritePacket(&packet.Text{
							Message: text.Colourf("<red>Unable to transfer: %s</red>", err.Error()),
						})
					}
					continue
				}
			case *packet.PlayerAction:
				if pk.ActionType == packet.PlayerActionDimensionChangeDone && s.transferring.CAS(true, false) {
					old := s.serverConn
					conn := s.tempServerConn

					pos := conn.GameData().PlayerPosition
					_ = s.conn.WritePacket(&packet.ChangeDimension{
						Dimension: packet.DimensionOverworld,
						Position:  pos,
					})

					_ = old.Close()

					s.serverConn = conn
					s.tempServerConn = nil

					s.updateTranslatorData(conn.GameData())

					// TODO: Set gamemode and stuff
					continue
				}
			}

			_ = s.ServerConn().WritePacket(pk)
		}
	}()

	go func() {
		defer func() {
			s.Close()
		}()
		for {
			conn := s.ServerConn()
			pk, err := conn.ReadPacket()
			if err != nil {
				log.Println(err)
				if conn != s.ServerConn() {
					continue
				}
				return
			}
			s.translatePacket(pk)

			switch pk := pk.(type) {
			case *packet.AvailableCommands:
				pk.Commands = append(pk.Commands, protocol.Command{
					Name:        "server",
					Description: "See the name of the server you are currently on",
				})

				var overloads []protocol.CommandOverload
				for _, g := range server.Groups() {
					var servers []string
					for _, s := range g.Servers() {
						servers = append(servers, s.Name())
					}
					fmt.Printf("Group %s has the servers %q\n", g.Name(), servers)
					overloads = append(overloads, protocol.CommandOverload{
						Parameters: []protocol.CommandParameter{
							{
								Name: "group",
								Type: protocol.CommandArgEnum | protocol.CommandArgValid,
								Enum: protocol.CommandEnum{
									Type:    g.Name() + "group",
									Options: []string{g.Name()},
								},
							},
							{
								Name: "server",
								Type: protocol.CommandArgEnum | protocol.CommandArgValid,
								Enum: protocol.CommandEnum{
									Type:    g.Name() + "server",
									Options: servers,
								},
							},
						},
					})
				}
				pk.Commands = append(pk.Commands, protocol.Command{
					Name:        "transfer",
					Description: "Transfer to another server on the proxy",
					/*Overloads: []protocol.CommandOverload{
						{
							Parameters: []protocol.CommandParameter{
								{
									Name: "server",
									Type: protocol.CommandArgEnum | protocol.CommandArgValid,
									Enum: protocol.CommandEnum{
										Type:    "servername",
										Options: servers,
									},
								},
								{
									Name: "group",
									Type: protocol.CommandArgEnum | protocol.CommandArgValid,
									Enum: protocol.CommandEnum{
										Type:    "groupname",
										Options: groups,
									},
								},
							},
						},
					},*/
					Overloads: overloads,
				})
			}

			_ = s.Conn().WritePacket(pk)
		}
	}()
}
