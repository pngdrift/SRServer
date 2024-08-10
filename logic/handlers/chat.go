package handlers

import (
	"log"
	"srserver/common"
	"srserver/logic/chat"
	"srserver/logic/clients"
	"srserver/logic/playermanager"
	pb "srserver/proto/out"
)

func SendChatMessage(pack *common.Pack, client *clients.Client) {
	channelType := pack.ReadString()
	message := &pb.ChatMessage{}
	pack.ReadProto(message)
	response := common.NewResponse(pack)
	response.SetMethod(common.OnNewChatMessage)
	response.WriteString(channelType)
	response.WriteProto(message)
	chat.AddMessage(message)
	Broadcast(response.ToByteBuff())
}

func SendPrivateChatMessage(pack *common.Pack, client *clients.Client) {
	message := &pb.ChatMessage{}
	pack.ReadProto(message)
	toUid := pack.ReadLong()
	log.Println("SendPrivateChatMessage", message, toUid)
	toPlayer := playermanager.FindPlayerById(toUid)
	if toPlayer == nil {
		return
	}
}
