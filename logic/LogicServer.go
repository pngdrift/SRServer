package logic

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"srserver/common"
	"srserver/logic/clients"
	"srserver/logic/db"
	"srserver/logic/handlers"
	"time"
)

type LogicServer struct {
}

func NewLogicServer() *LogicServer {
	return &LogicServer{}
}

func (s *LogicServer) Start() error {
	db.InitDb()
	ln, err := net.Listen("tcp", ":8992")
	if err != nil {
		log.Fatalf("Failed to start content Logic: %v", err)
	}
	log.Println("Logic server is listening on port 8992")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Logic server err: %v", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *LogicServer) Shutdown() {
	log.Println("Shutdown server")
	seconds := []int32{30, 20, 10, 3, 0}
	for _, second := range seconds {
		time.AfterFunc(time.Duration(seconds[0]-second)*time.Second, func() {
			if second == 0 {
				os.Exit(0)
			}
			pack := common.NewPack()
			pack.SetMethod(common.OnShutdown)
			pack.WriteInt(second)
			handlers.Broadcast(pack.ToByteBuff())
		})
	}
}
func (s *LogicServer) handleConnection(conn net.Conn) {
	clients.Clients[conn] = &clients.Client{
		Conn:   conn,
		Player: nil,
	}
	defer func() {
		delete(clients.Clients, conn)
		conn.Close()
	}()
	reader := bufio.NewReader(conn)
	packetLengthBuffer := make([]byte, 4)
	for {
		_, err := io.ReadFull(reader, packetLengthBuffer)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading packet length:", err)
			}
			return
		}

		packetLength := binary.BigEndian.Uint32(packetLengthBuffer)

		packetBuffer := make([]byte, packetLength)
		_, err = io.ReadFull(reader, packetBuffer)
		if err != nil {
			log.Println("Error reading packet:", err)
			return
		}

		pack := common.NewPackInstance(packetBuffer)
		handlers.ProcessPacket(pack, clients.Clients[conn])
	}
}
