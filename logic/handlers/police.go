package handlers

import (
	"math/rand/v2"
	"srserver/common"
	"srserver/conf"
	"srserver/logic/clients"
	pb "srserver/proto/out"
)

func RegisterCar(pack *common.Pack, client *clients.Client) {
	carId, region := pack.ReadLong(), pack.ReadInt()

	carNumber := &pb.CarNumber{
		CarId:     carId,
		IsTransit: false,
		Number:    0,
		Region:    region,
	}

	userCar := client.Player.FindUserCarById(carId)
	if userCar == nil {
		return
	}
	userCar.Number = carNumber
	client.Player.Save()
	response := common.NewResponse(pack)
	response.WriteProto(carNumber)
	client.Write(response.ToByteBuff())
}

func GetNumberShop(pack *common.Pack, client *clients.Client) {
	region := pack.ReadInt()
	response := common.NewResponse(pack)
	response.WriteInt(region)
	len := rand.Int32N(100)
	response.WriteInt(len)
	for i := int32(0); i < len; i++ {
		response.WriteInt(i)
	}
	client.Write(response.ToByteBuff())
}

func BuyNumber(pack *common.Pack, client *clients.Client) {
	carId := pack.ReadLong()
	number := pack.ReadInt()
	region := pack.ReadInt()

	userCar := client.Player.FindUserCarById(carId)
	if userCar == nil {
		return
	}

	carNumber := &pb.CarNumber{
		CarId:     carId,
		IsTransit: false,
		Number:    number,
		Region:    region,
	}
	userCar.Number = carNumber
	client.Player.Withdraw(conf.SHOP_NUMBERS_COST)

	response := common.NewResponse(pack)
	response.WriteProto(carNumber)
	client.Write(response.ToByteBuff())
}

func SwapNumbers(pack *common.Pack, client *clients.Client) {
	fromCarId := pack.ReadLong()
	toCarId := pack.ReadLong()
	fromCar := client.Player.FindUserCarById(fromCarId)
	toCar := client.Player.FindUserCarById(toCarId)
	if fromCar == nil || toCar == nil {
		return
	}
	number1 := fromCar.Number
	number2 := toCar.Number

	number2.CarId = fromCar.Id
	fromCar.Number = number2

	number1.CarId = toCar.Id
	toCar.Number = number1

	client.Player.Save()
}
