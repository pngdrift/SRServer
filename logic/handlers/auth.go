package handlers

import (
	"log"
	"srserver/common"
	"srserver/conf"
	"srserver/logic/clients"
	"srserver/logic/db"
	"srserver/logic/playermanager"
	pb "srserver/proto/out"
)

func CheckVersion(pack *common.Pack, client *clients.Client) {
	clientVersion := pack.ReadInt()
	response := common.NewResponse(pack)
	if clientVersion < conf.MIN_CLIENT_VERSION || clientVersion > conf.MAX_CLIENT_VERSION {
		response.SetError(true)
	} else {
		response.WriteLong(db.DatabaseChecksum) //todo check checksum type
	}
	client.Write(response.ToByteBuff())
}

func Load(pack *common.Pack, client *clients.Client) {
	response := common.NewResponse(pack)
	response.WriteBytes(db.DatabaseBuff)
	client.Write(response.ToByteBuff())
}

func CheckBan(pack *common.Pack, client *clients.Client) {
	socialType, id := pack.ReadString(), pack.ReadString()
	log.Println("CheckBan", socialType, id)
	response := common.NewResponse(pack)
	if _, exists := playermanager.BannedPlayers[id+socialType]; exists {
		response.SetError(true)
		response.WriteLong(0) //uid
		response.WriteString("You are banned.")
		response.WriteProto(&pb.GameException{
			ErrorMessage:     "Ban",
			ErrorDescription: "BANNED",
			Error:            0,
			ErrorLevel:       pb.ErrorLevel_INFO,
		})
	}

	client.Write(response.ToByteBuff())
}

func Init(pack *common.Pack, client *clients.Client) {
	socialType := pack.ReadString()
	id := pack.ReadString()
	timeZone := pack.ReadString()
	language := pack.ReadString()
	currentTime := pack.ReadLong()
	log.Printf("Init player: socialType: %s id: %s timezone: %s lang: %s current time: %d", socialType, id, timeZone, language, currentTime)
	client.Player = playermanager.Init(id+socialType, client)
	user := client.Player.User
	user.World.ZoneId = timeZone
	response := common.NewResponse(pack)
	response.WriteProto(user)
	response.WriteLong(0) //detla time

	client.Write(response.ToByteBuff())
}
