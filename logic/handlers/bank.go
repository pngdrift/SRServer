package handlers

import (
	"log"
	"srserver/common"
	"srserver/logic/bank"
	"srserver/logic/clients"
	pb "srserver/proto/out"
	"time"
)

func GetPurchaseList(pack *common.Pack, client *clients.Client) {
	response := common.NewResponse(pack)
	response.WriteString("qiwi")
	response.WriteString("tele2")
	client.Write(pack.ToByteBuff())
}

func Exchange(pack *common.Pack, client *clients.Client) {
	itemId := pack.ReadString()
	exchangeItem := bank.FindExchangeItem(itemId)
	if exchangeItem == nil ||
		!client.Player.CheckMoney(exchangeItem.FromMoney) ||
		client.Player.User.TimersAndCounters.ExchangeCount == 3 {
		return
	}

	client.Player.User.TimersAndCounters.ExchangeCount++
	client.Player.Withdraw(exchangeItem.FromMoney)
	client.Player.Deposit(exchangeItem.ToMoney)

	response := common.NewResponse(pack)
	response.WriteProto(exchangeItem)
	client.Write(response.ToByteBuff())
}

func GetBank(pack *common.Pack, client *clients.Client) {
	if pack.IsHasBytes() {
		item := &pb.AndroidBankItem{}
		pack.ReadProto(item)
	}
	response := common.NewResponse(pack)
	response.WriteProto(bank.GetBank())
	client.Write(response.ToByteBuff())
}

func LoadWallet(pack *common.Pack, client *clients.Client) {
	uid := pack.ReadLong()
	limit := pack.ReadInt()
	var fromTransaction int64
	if pack.IsHasLong() {
		fromTransaction = pack.ReadLong()
	}
	log.Println("LoadWallet", uid, limit, fromTransaction)
	wallet := &pb.Wallet{
		Id: uid,
		CommonLog: []*pb.Transaction{{
			Id:          1,
			Type:        pb.TransactionType_BUY,
			MoneyBefore: &pb.Money{Gold: 200, Money: 500},
			MoneyAfter:  &pb.Money{Gold: 300},
			Method:      0,
			Description: "Test",
			Time:        time.Now().UnixMilli(),
		}},
		BankLog: []*pb.Transaction{},
	}
	response := common.NewResponse(pack)
	response.WriteProto(wallet)
	client.Write(response.ToByteBuff())
}
