package handlers

import (
	"log"
	"srserver/common"
	"srserver/conf"
	"srserver/logic/clients"
	"srserver/logic/db"
	"srserver/logic/playermanager"
	"srserver/logic/utils"
	pb "srserver/proto/out"
)

func BuyCar(pack *common.Pack, client *clients.Client) {
	baseId := pack.ReadInt()
	for _, userCar := range client.Player.User.Garage.Cars {
		if userCar.BaseId == baseId {
			return
		}
	}
	car := db.FindCar(baseId)
	if car != nil {
		response := common.NewResponse(pack)
		if client.Player.CheckMoney(car.GetPrice()) {
			response.WriteProto(client.Player.BuyCar(car))
			client.Write(response.ToByteBuff())
		}
	}
}

func BuyItem(pack *common.Pack, client *clients.Client) {
	baseId := pack.ReadInt()
	count := pack.ReadInt()
	itemType := pb.ItemType(pb.ItemType_value[pack.ReadString()])
	var cost *pb.Money
	switch itemType {
	case pb.ItemType_BLUEPRINT:
		blueprint := db.FindBlueprint(baseId)
		if blueprint != nil {
			cost = utils.MultiplyMoney(blueprint.GetBase().Price, count)
		}
	case pb.ItemType_TOOLS:
		tool := db.FindTools(baseId)
		if tool != nil {
			cost = utils.MultiplyMoney(tool.GetBase().Price, count)
		}
	}
	player := client.Player
	if cost == nil || count <= 0 || !player.CheckMoney(cost) {
		return
	}
	item := &pb.Item{
		Id:     utils.MakeId(), //todo maybe id == baseId or smth like
		Count:  count,
		Type:   itemType,
		BaseId: baseId,
	}
	player.Withdraw(cost)
	response := common.NewResponse(pack)
	response.WriteProto(item)
	client.Write(response.ToByteBuff())

	inventory := player.User.Inventory.Items
	var itemExists bool
	for _, invItem := range inventory {
		if invItem.BaseId == baseId && invItem.Type == itemType {
			itemExists = true
			invItem.Count += count
			break
		}
	}
	if !itemExists {
		item.Id = int64(baseId)
		player.User.Inventory.Items = append(inventory, item)
	}
	player.Save()
}

func SellItem(pack *common.Pack, client *clients.Client) {
	itemId := pack.ReadLong()
	count := pack.ReadInt()
	inventory := client.Player.User.Inventory.Items
	for i, item := range inventory {
		if item.Id == itemId {
			if item.Count >= count {
				item.Count -= count
				if item.Count == 0 {
					client.Player.User.Inventory.Items = append(inventory[:i], inventory[i+1:]...)
				}
				var base *pb.BaseItem
				if item.Type == pb.ItemType_BLUEPRINT {
					base = db.FindBlueprint(item.BaseId).GetBase()
				} else {
					base = db.FindTools(item.BaseId).GetBase()
				}
				client.Player.Deposit(base.Price) //todo divide price
			}
			break
		}
	}
}

func SellUpgrades(pack *common.Pack, client *clients.Client) {
	player := client.Player
	for pack.IsHasLong() {
		upgradeId := pack.ReadLong()
		for _, item := range player.User.Inventory.Upgrades {
			if item.Id == upgradeId {
				//todo find base item
				player.Deposit(&pb.Money{})
			}
		}
	}
}

func PaintCar(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	player := client.Player
	userCar := player.FindUserCarById(carId)
	if userCar == nil {
		return
	}
	paintCommands := &pb.PaintCommands{}
	pack.ReadProto(paintCommands)
	for _, command := range paintCommands.Commands {
		if len(command.GetIntArgs()) == 0 {
			if command.Type == pb.PaintCmdType_PAINT_DISK {
				userCar.Paint.IsDiskPainted = false
			} else if command.Type == pb.PaintCmdType_PAINT_DISK_FRONT {
				userCar.Paint.IsDiskPaintedFront = false
			}
			continue
		}
		switch command.Type {
		case pb.PaintCmdType_PAINT_CHASSIS:
			userCar.Paint.CarColor = command.IntArgs[0]
			buyColor(command.IntArgs[0], player)
		case pb.PaintCmdType_PAINT_FRONT_BUMPER:
			userCar.Paint.FrontBumperColor = command.IntArgs[0]
			buyColor(command.IntArgs[0], player)
		case pb.PaintCmdType_PAINT_CENTER_BUMPER:
			userCar.Paint.CenterBumperColor = command.IntArgs[0]
			buyColor(command.IntArgs[0], player)
		case pb.PaintCmdType_PAINT_REAR_BUMPER:
			userCar.Paint.RearBumperColor = command.IntArgs[0]
			buyColor(command.IntArgs[0], player)
		case pb.PaintCmdType_INSTALL_TINT:
			userCar.Paint.TintingColor = command.IntArgs[0]
			buyColor(command.IntArgs[0], player)
		case pb.PaintCmdType_PAINT_DISK:
			userCar.Paint.DiskColor = command.IntArgs[0]
			userCar.Paint.IsDiskPainted = true
			buyColor(command.IntArgs[0], player)
		case pb.PaintCmdType_PAINT_DISK_FRONT:
			userCar.Paint.DiskColorFront = command.IntArgs[0]
			userCar.Paint.IsDiskPaintedFront = true
			buyColor(command.IntArgs[0], player)
		case pb.PaintCmdType_PAINT_RIM:
			userCar.Paint.RimColor = command.IntArgs[0]
			buyColor(command.IntArgs[0], player)
		case pb.PaintCmdType_PAINT_RIM_FRONT:
			userCar.Paint.RimColorFront = command.IntArgs[0]
			buyColor(command.IntArgs[0], player)
		case pb.PaintCmdType_ADD_DEACL:
			baseDecal := db.FindDecal(command.IntArgs[1])
			if baseDecal != nil {
				userCar.Paint.Decals = append(userCar.Paint.Decals, &pb.Decal{
					Id:     command.IntArgs[0],
					BaseId: command.IntArgs[1],
				})
				player.Withdraw(baseDecal.GetPrice())
			}
		case pb.PaintCmdType_ADD_USER_DECAL: // Not working fully in client
			baseDecal := db.FindDecal(10000)
			log.Println("ADD_USER_DECAL", command)
			userCar.Paint.Decals = append(userCar.Paint.Decals, &pb.Decal{
				Id:        command.IntArgs[0],
				BaseId:    10000,
				UserDecal: true,
				FileName:  "testUserDecal",
			})
			player.Withdraw(baseDecal.GetPrice())
		case pb.PaintCmdType_REMOVE_DECAL:
			for i, decal := range userCar.Paint.Decals {
				if decal.GetId() == command.IntArgs[0] {
					userCar.Paint.Decals = append(userCar.Paint.Decals[:i], userCar.Paint.Decals[i+1:]...)
					break
				}
			}
		case pb.PaintCmdType_PAINT_DECAL:
			for _, decal := range userCar.Paint.Decals {
				if decal.GetId() == command.IntArgs[0] {
					decal.Color = command.IntArgs[1]
					break
				}
			}
		case pb.PaintCmdType_UPDATE_DECAL:
			if len(command.FloatArgs) < 4 {
				continue
			}
			for _, decal := range userCar.Paint.Decals {
				if decal.GetId() == command.IntArgs[0] {
					decal.X = command.FloatArgs[0]
					decal.Y = command.FloatArgs[1]
					decal.Scale = command.FloatArgs[2]
					decal.Rotation = command.FloatArgs[3]
					break
				}
			}
		}
		player.Save()
		client.Write(pack.ToByteBuff())
	}
}

func buyColor(baseId int32, player *playermanager.Player) {
	color := db.FindCar(baseId)
	if color != nil {
		player.Withdraw(color.GetPrice())
	}
}

func GetUserDecal(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	fileName := pack.ReadString()
	log.Println("GetUserDecal", carId, fileName)
	// response := common.NewResponse(pack)
	// response.WriteBytes(testDecal)
	// client.Write(response.ToByteBuff())
}

func Refuel(pack *common.Pack, client *clients.Client) {
	tc := client.Player.User.TimersAndCounters
	if tc.RefuelCount < 5 {
		client.Player.Withdraw(conf.REFUEL_COST[tc.RefuelCount])
		client.Player.AddFuel(conf.REFUEL_VALUE)
		tc.RefuelCount++
	}
}

func SellCar(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	garage := client.Player.User.Garage
	for i, userCar := range garage.Cars {
		if userCar.Id == carId {
			baseCar := db.FindCar(userCar.BaseId)
			garage.Cars = append(garage.Cars[:i], garage.Cars[i+1:]...)
			client.Player.Deposit(baseCar.GetPrice()) //todo
			break
		}
	}
}

func SelectCar(pack *common.Pack, client *clients.Client) {
	client.Player.User.Garage.CurrentId = pack.ReadLong()
}

func WashCar(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	player := client.Player
	userCar := player.FindUserCarById(carId)
	if userCar == nil {
		return
	}
	userCar.Dirtiness = 0
	player.Withdraw(&pb.Money{Money: 800})
}
