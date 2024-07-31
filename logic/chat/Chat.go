package chat

import (
	pb "srserver/proto/out"
	"time"
)

var publicMessages = []*pb.ChatMessage{
	{
		Id:      0,
		Time:    0,
		FromUid: 0,
		Message: "Plz not insults here",
		UserInfo: &pb.UserInfo{
			Name:   "System",
			Avatar: "",
			Id:     0,
			Lang:   "",
			Type:   0,
		},
		ToUid:      0,
		ToUserInfo: &pb.UserInfo{},
		Registred:  false,
	},
}

func AddMessage(message *pb.ChatMessage) {
	message.FromUid = message.UserInfo.GetId()
	message.Time = time.Now().UnixMilli()
	message.Id = int64(len(publicMessages))
	publicMessages = append(publicMessages, message)
}

func GetInitData() *pb.Chat {
	return &pb.Chat{
		Rooms: []*pb.ChatRoom{
			{
				Id:       pb.ChatRoomType_PUBLIC.String(),
				IsBanned: false,
				IsLoaded: true,
				IsLocked: false,
				Type:     pb.ChatRoomType_PUBLIC,
				Messages: publicMessages,
			},
			{
				Id:       pb.ChatRoomType_TAXI.String(),
				IsBanned: false,
				IsLoaded: true,
				IsLocked: false,
				Type:     pb.ChatRoomType_TAXI,
				Messages: []*pb.ChatMessage{},
			},
			{
				Id:       pb.ChatRoomType_PRIVATE.String(),
				IsBanned: false,
				IsLoaded: true,
				IsLocked: false,
				Type:     pb.ChatRoomType_PRIVATE,
				Messages: []*pb.ChatMessage{},
			},
		},
	}
}
