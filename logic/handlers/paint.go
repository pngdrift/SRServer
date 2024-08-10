package handlers

import (
	"log"
	"srserver/common"
	"srserver/logic/clients"
	"srserver/logic/utils"
	pb "srserver/proto/out"
	"strconv"
)

func SaveCurrentPaint(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	player := client.Player
	userCar := player.FindUserCarById(carId)
	if userCar == nil {
		return
	}
	//todo check if livery already been saved

	paintId := utils.MakeId()
	paintItem := &pb.PaintItem{
		Id:        paintId,
		BaseCarId: userCar.BaseId,
		Name:      "Livery #" + strconv.FormatInt(paintId, 10),
		Paint:     userCar.Paint,
		IsShared:  false,
	}
	player.User.Paints.Items = append(player.User.Paints.Items, paintItem)
	player.Save()
	response := common.NewResponse(pack)
	response.WriteLong(carId)
	response.WriteProto(paintItem)
	client.Write(response.ToByteBuff())
}

func RenamePaintItem(pack *common.Pack, client *clients.Client) {
	paintId := pack.ReadLong()
	name := pack.ReadString()
	player := client.Player
	paintItem := player.FindPaintItemById(paintId)
	if paintItem == nil {
		return
	}
	paintItem.Name = name
	player.Save()
}

func ApplyPaint(pack *common.Pack, client *clients.Client) {
	paintId := pack.ReadLong()
	carId := pack.ReadLong()
	player := client.Player
	userCar := player.FindUserCarById(carId)
	if userCar == nil {
		return
	}
	paintItem := player.FindPaintItemById(paintId)
	if paintItem == nil {
		return
	}
	userCar.Paint = paintItem.Paint
	player.Save()
}

func DeletePaint(pack *common.Pack, client *clients.Client) {
	paintId := pack.ReadLong()
	player := client.Player
	userPaints := player.User.Paints
	for i, item := range userPaints.Items {
		if item.Id == paintId {
			userPaints.Items = append(userPaints.Items[:i], userPaints.Items[i+1:]...)
			player.Save()
			break
		}
	}
}

func SharePaint(pack *common.Pack, client *clients.Client) {
	paintId := pack.ReadLong()
	player := client.Player
	paintItem := player.FindPaintItemById(paintId)
	if paintItem == nil {
		return
	}
	paintItem.IsShared = true
	//todo store in db
	player.Save()
}

func ObtainPaint(pack *common.Pack, client *clients.Client) {
	paintId := pack.ReadLong()

	response := common.NewResponse(pack)
	response.SetError(true)
	response.WriteProto(&pb.GameException{
		ErrorMessage: "Not implemented",
	})
	client.Write(response.ToByteBuff())
	log.Println("ObtainPaint", paintId)
}

func UnSharePaint(pack *common.Pack, client *clients.Client) {
	paintId := pack.ReadLong()
	player := client.Player
	paintItem := player.FindPaintItemById(paintId)
	if paintItem == nil {
		return
	}
	paintItem.IsShared = false
	player.Save()
}
