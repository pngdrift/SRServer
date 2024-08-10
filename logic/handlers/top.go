package handlers

import (
	"srserver/common"
	"srserver/logic/clients"
	"srserver/logic/top"
)

func UpdateTop(pack *common.Pack, client *clients.Client) {
	response := common.NewResponse(pack)
	response.WriteProto(top.GetInitData(true))
	client.Write(response.ToByteBuff())
}
