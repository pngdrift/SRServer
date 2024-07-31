package clients

import (
	"net"
	"srserver/logic/playermanager"
)

type Client struct {
	net.Conn
	Player *playermanager.Player
}

var Clients map[net.Conn]*Client = make(map[net.Conn]*Client)
