package handlers

import (
	"log"
	"math/rand"
	"srserver/common"
	"srserver/conf"
	"srserver/logic/clients"
	"srserver/logic/utils"
	pb "srserver/proto/out"
)

func UpgradeEngine(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	slotName := pack.ReadString()
	player := client.Player
	userCar := player.FindUserCarById(carId)
	if userCar == nil {
		return
	}
	var engineUpgradeType pb.EngineUpgradeType = pb.EngineUpgradeType(pb.EngineUpgradeType_value[slotName])

	oldLevel := userCar.GetEngineUpgrades()[engineUpgradeType.Number()].Level
	if oldLevel >= 10 {
		return
	}
	newLevel := oldLevel + 1
	cost := conf.UPGRADE_COST_MAP[engineUpgradeType][newLevel]
	if player.CheckMoney(cost) {
		userCar.GetEngineUpgrades()[engineUpgradeType.Number()].Level = newLevel
		player.Withdraw(cost)
	}
}

func BuyUpgrade(pack *common.Pack, client *clients.Client) {
	baseId := pack.ReadInt()
	upgradeType := pb.UpgradeType(pb.UpgradeType_value[pack.ReadString()])
	carUpgrade := &pb.CarUpgrade{
		Id:       utils.MakeId(),
		CarId:    0,
		Current:  0,
		BaseId:   baseId,
		Type:     upgradeType,
		IsPacked: true,
		Grade:    pb.UpgradeGrade_WHITE,
	}
	inventory := client.Player.User.Inventory
	inventory.Upgrades = append(inventory.Upgrades, carUpgrade)
	client.Player.Withdraw(&pb.Money{})
	response := common.NewResponse(pack)
	response.WriteProto(carUpgrade)
	client.Write(response.ToByteBuff())
}

func InstallUpgrade(pack *common.Pack, client *clients.Client) {
	carId, upgradeId := pack.ReadLong(), pack.ReadLong()
	slotName := pack.ReadString()
	player := client.Player
	userCar := player.FindUserCarById(carId)
	if userCar == nil {
		return
	}
	inventory := player.User.Inventory
	for i, upgrade := range inventory.Upgrades {
		if upgrade.Id == upgradeId && (upgrade.CarId == carId || upgrade.IsPacked) {
			upgrade.IsPacked = false
			upgrade.CarId = carId
			upgradeSlot := &pb.UpgradeSlot{
				CarId:   carId,
				Type:    utils.GetUpgradeType(slotName),
				Upgrade: upgrade,
			}
			utils.InstallUpgrade(userCar, upgradeSlot, slotName)
			inventory.Upgrades = append(inventory.Upgrades[:i], inventory.Upgrades[i+1:]...)
			player.Save()
			break
		}
	}
}

func UninstallUpgrade(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	slotName := pack.ReadString()
	player := client.Player
	userCar := player.FindUserCarById(carId)
	if userCar == nil {
		return
	}
	upgradeSlot := utils.GetUpgradeSlot(userCar, slotName)
	if upgradeSlot == nil {
		return
	}
	player.User.Inventory.Upgrades = append(player.User.Inventory.Upgrades, upgradeSlot.Upgrade)
	utils.InstallUpgrade(userCar, nil, slotName)
	player.Save()
}

func ConfigUpgrades(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	frontSpring, rearSuspension := pack.ReadFloat(), pack.ReadFloat()
	rearSpring, frontSuspension := pack.ReadFloat(), pack.ReadFloat()
	clirence := pack.ReadFloat()
	//todo check allowed values for installed springs

	userCar := client.Player.FindUserCarById(carId)
	if userCar == nil {
		return
	}
	setting := &pb.Setting{
		Id:              rand.Int31(),
		IsActive:        true,
		FrontSpring:     frontSpring,
		FrontSuspension: frontSuspension,
		RearSpring:      rearSpring,
		RearSuspension:  rearSuspension,
		Clirence:        clirence,
	}

	userCar.Settings.Settings = append(userCar.Settings.Settings, setting)
}

func ConfigMuffler(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	offsetX := pack.ReadFloat()
	offsetY := pack.ReadFloat()
	userCar := client.Player.FindUserCarById(carId)
	if userCar == nil ||
		offsetX < -0.05 || offsetY < -0.05 ||
		offsetX > 0.1 || offsetY > 0.15 {
		return
	}
	userCar.Settings.MufflerOffsetX = offsetX
	userCar.Settings.MufflerOffsetY = offsetY
	client.Player.Save()
}

func ConfigShiftLamps(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	yellowZone := pack.ReadInt()
	greenZone := pack.ReadInt()
	redZone := pack.ReadInt()
	userCar := client.Player.FindUserCarById(carId)
	if userCar == nil {
		return
	}
	userCar.Settings.YellowZoneRpm = yellowZone
	userCar.Settings.GreenZoneRpm = greenZone
	userCar.Settings.RedZoneRpm = redZone
	client.Player.Save()
}

func CraftUpgrade(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	upgradeId := pack.ReadLong()
	blueprintId := pack.ReadLong()
	toolsId := pack.ReadLong()
	log.Println("CraftUpgrade", carId, upgradeId, blueprintId, toolsId)
	craftResult := &pb.CraftResult{
		Success:     false,
		CarId:       carId,
		UpgradeId:   upgradeId,
		BlueprintId: blueprintId,
		ToolsId:     toolsId,
		Upgrade: &pb.CarUpgrade{
			Id:       utils.MakeId(),
			CarId:    carId,
			Current:  0,
			BaseId:   0, //todo search item
			Type:     0,
			IsPacked: false,
			Grade:    0,
		},
	}
	response := common.NewResponse(pack)
	response.WriteProto(craftResult)
	client.Write(response.ToByteBuff())
}

func UpdateWheelPosition(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	frontX := pack.ReadFloat()
	frontY := pack.ReadFloat()
	rearX := pack.ReadFloat()
	rearY := pack.ReadFloat()
	player := client.Player
	userCar := player.FindUserCarById(carId)
	if userCar == nil {
		return
	}
	userCar.FrontWheelX = frontX
	userCar.FrontWheelY = frontY
	userCar.RearWheelX = rearX
	userCar.RearWheelY = rearY
	player.Save()
}
