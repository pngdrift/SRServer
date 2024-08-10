package handlers

import (
	"srserver/common"
	"srserver/conf"
	"srserver/logic/clients"
	"srserver/logic/utils"
	pb "srserver/proto/out"
)

func SaveDynoTest(pack *common.Pack, client *clients.Client) {
	test := &pb.DynoTest{}
	pack.ReadProto(test)
	client.Player.User.Dyno.CurrentTest = test
	client.Player.User.Dyno.CurrentTestId = utils.MakeId()
	client.Player.Save()
	response := common.NewResponse(pack)
	response.WriteProto(client.Player.User.Dyno)
	client.Write(response.ToByteBuff())
}

func StartDynoTest(pack *common.Pack, client *clients.Client) {
	response := common.NewResponse(pack)
	client.Write(response.ToByteBuff())
	client.Player.Withdraw(conf.MONEY_FOR_DYNO_TEST)
}
