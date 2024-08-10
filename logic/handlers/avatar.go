package handlers

import (
	"log"
	"srserver/common"
	"srserver/logic/clients"
	"srserver/logic/playermanager"
	"srserver/logic/utils"
	pb "srserver/proto/out"
	"strconv"
)

func UpdateUserInfo(pack *common.Pack, client *clients.Client) {
	userInfo := &pb.UserInfo{}
	pack.ReadProto(userInfo)
	client.Player.User.Info = userInfo //todo check correctly
	client.Player.Save()
}

func GetAvatar(pack *common.Pack, client *clients.Client) {
	userId := pack.ReadLong()
	player := playermanager.FindPlayerById(userId)
	response := common.NewResponse(pack)
	var url string
	if player != nil && len(player.User.Info.Avatar) > 10 {
		url = player.User.Info.Avatar
	} else {
		url = "https://robohash.org/" + strconv.FormatInt(userId, 16)
	}
	avatar, err := utils.LoadAvatar(url)
	if err != nil {
		log.Println("GetAvatar err", err)
	}
	response.WriteBytes(avatar)
	client.Write(response.ToByteBuff())
}
