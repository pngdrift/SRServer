package handlers

import (
	"srserver/common"
	"srserver/logic/clients"
)

func RegisterInTournament(pack *common.Pack, client *clients.Client) {
	//tournamentId := pack.ReadLong()
	response := common.NewResponse(pack)
	response.SetError(false)
	client.Write(response.ToByteBuff())
}

func GetTournamentEnemy(pack *common.Pack, client *clients.Client) {
	//tournamentId := pack.ReadLong()
}
