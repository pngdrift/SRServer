package handlers

import (
	"log"
	"srserver/common"
	"srserver/logic/clients"
	pb "srserver/proto/out"
)

func CompleteQuest(pack *common.Pack, client *clients.Client) {
	id := pack.ReadInt()
	log.Print("CompleteQuest ", id)
	response := common.NewResponse(pack)
	award := &pb.QuestAward{
		Exp: 100,
		Money: &pb.Money{
			Money:            110,
			Gold:             0,
			TournamentPoints: 0,
			TopPoints:        0,
			UpgradePoints:    0,
		},
	}
	response.WriteProto(award)
	client.Write(response.ToByteBuff())
}
