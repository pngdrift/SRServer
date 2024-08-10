package handlers

import (
	"log"
	"srserver/common"
	"srserver/logic/clients"
	"srserver/logic/mail"
	"srserver/logic/playermanager"
	pb "srserver/proto/out"
)

func AddGameBan(pack *common.Pack, client *clients.Client) {
	parentId, uid := pack.ReadLong(), pack.ReadLong()
	// idk reason = pack.ReadInt()
	if parentId != client.Player.User.Id {
		return
	}
	log.Println(parentId, uid)
	player := playermanager.FindPlayerById(uid)
	if player != nil {
		log.Println("AddGameBan")
		playermanager.BannedPlayers[player.Uid] = true
		if player.Conn != nil {
			player.Conn.Close()
		}
	}
	response := common.NewResponse(pack)
	client.Write(response.ToByteBuff())
}

func RemoveGameBan(pack *common.Pack, client *clients.Client) {
	parentId, uid := pack.ReadLong(), pack.ReadLong()
	if parentId != client.Player.User.Id {
		return
	}
	player := playermanager.FindPlayerById(uid)
	if player != nil {
		log.Println("RemoveGameBan")
		delete(playermanager.BannedPlayers, player.Uid)
	}
	response := common.NewResponse(pack)
	client.Write(response.ToByteBuff())
}

func SendSystemMail(pack *common.Pack, client *clients.Client) {
	parentId, uid := pack.ReadLong(), pack.ReadLong()
	if parentId != client.Player.User.Id {
		return
	}
	money := &pb.Money{}
	pack.ReadProto(money)
	player := playermanager.FindPlayerById(uid)
	response := common.NewResponse(pack)
	if player != nil {
		mail.SendSystemMail(*player, money)
	} else {
		response.SetError(true)
		response.WriteProto(&pb.GameException{
			ErrorMessage:     "uid incorrect",
			ErrorDescription: "test",
			Error:            0,
			ErrorLevel:       pb.ErrorLevel_INFO,
		})
	}
	client.Write(response.ToByteBuff())
}
