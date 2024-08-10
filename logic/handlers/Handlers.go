package handlers

import (
	"log"
	"srserver/common"
	"srserver/logic/clients"
)

type PacketHandlerFunc func(*common.Pack, *clients.Client)

var Handlers = map[int32]PacketHandlerFunc{
	common.CheckVersion: CheckVersion,
	common.Load:         Load,
	common.Init:         Init,
	common.CheckBan:     CheckBan,

	common.RegisterCar:    RegisterCar,
	common.BuyNumber:      BuyNumber,
	common.GetNumbersShop: GetNumberShop,
	common.SwapNumbers:    SwapNumbers,

	common.CompleteQuest: CompleteQuest,

	common.SaveDynoTest:  SaveDynoTest,
	common.StartDynoTest: StartDynoTest,

	common.BuyCar:       BuyCar,
	common.BuyItem:      BuyItem,
	common.SellItem:     SellItem,
	common.SellUpgrades: SellUpgrades,
	common.PaintCar:     PaintCar,
	common.GetUserDecal: GetUserDecal,
	common.Refuel:       Refuel,
	common.SellCar:      SellCar,
	common.SelectCar:    SelectCar,
	common.WashCar:      WashCar,

	common.UpgradeEngine:       UpgradeEngine,
	common.BuyUpgrade:          BuyUpgrade,
	common.InstallUpgrade:      InstallUpgrade,
	common.UninstallUpgrade:    UninstallUpgrade,
	common.ConfigUpgrades:      ConfigUpgrades,
	common.ConfigMuffler:       ConfigMuffler,
	common.ConfigShiftLamps:    ConfigShiftLamps,
	common.CraftUpgrade:        CraftUpgrade,
	common.UpdateWheelPosition: UpdateWheelPosition,

	common.Exchange:        Exchange,
	common.GetBank:         GetBank,
	common.GetPurchaseList: GetPurchaseList,
	common.LoadWallet:      LoadWallet,

	common.UpdateUserInfo: UpdateUserInfo,
	common.GetAvatar:      GetAvatar,

	common.SendChatMessage:        SendChatMessage,
	common.SendPrivateChatMessage: SendPrivateChatMessage,

	common.UpdateEnemies:       UpdateEnemies,
	common.GetRaceAward:        GetRaceAward,
	common.CreateRace:          CreateRace,
	common.BrakeRace:           BrakeRace,
	common.UpdatePointsEnemies: UpdatePointsEnemies,

	common.UpdateTop: UpdateTop,

	common.ReadMail:          ReadMail,
	common.ReadAllMails:      ReadAllMails,
	common.DeleteMail:        DeleteMail,
	common.DeleteReadedMails: DeleteReadedMails,

	common.RegisterInTournament: RegisterInTournament,
	common.GetTournamentEnemy:   GetTournamentEnemy,

	common.SaveCurrentPaint: SaveCurrentPaint,
	common.RenamePaintItem:  RenamePaintItem,
	common.ApplyPaint:       ApplyPaint,
	common.DeletePaint:      DeletePaint,
	common.SharePaint:       SharePaint,
	common.ObtainPaint:      ObtainPaint,
	common.UnSharePaint:     UnSharePaint,

	common.AddGameBan:     AddGameBan,
	common.RemoveGameBan:  RemoveGameBan,
	common.SendSystemMail: SendSystemMail,
}

func ProcessPacket(pack *common.Pack, client *clients.Client) {
	//unauthorised client cant access these methods
	if pack.GetMethod() >= 200 && client.Player == nil {
		return
	}

	if handler, exists := Handlers[pack.GetMethod()]; exists {
		handler(pack, client)
	} else {
		log.Printf("No handler for method %d", pack.GetMethod())
	}
}

func Broadcast(b []byte) {
	for conn, client := range clients.Clients {
		if client.Player != nil {
			conn.Write(b)
		}
	}
}
