package bank

import pb "srserver/proto/out"

var bank = &pb.Bank{
	ExchangeItems: []*pb.ExchangeItem{
		{
			FromMoney: &pb.Money{
				Money: 10000,
			},
			ToMoney: &pb.Money{
				Gold: 10,
			},
			Order:  1,
			ItemId: "exchange1",
		},
		{
			FromMoney: &pb.Money{
				Gold:          7,
				UpgradePoints: 90,
			},
			ToMoney: &pb.Money{
				Money: 5000,
			},
			Order:  2,
			ItemId: "exchange2",
		},
		{
			FromMoney: &pb.Money{
				Gold: 100,
			},
			ToMoney: &pb.Money{
				TournamentPoints: 200,
			},
			Order:  3,
			ItemId: "exchange3",
		},
	},
	Items: []*pb.BankItem{
		/*{
			ItemId:   "bank_item",
			Money:    1000,
			Gold:     1000,
			Bonus:    20,
			Fuel:     50,
			LootId:   2,
			CarId:    1,
			Price:    "100 rub",
			Special:  false,
			Date:     "15.10.2023",
			Duration: 2100000,
			Revenue:  1.533,
		},*/
	},
}

func GetBank() *pb.Bank {
	return bank
}

func FindExchangeItem(itemId string) *pb.ExchangeItem {
	for _, item := range bank.GetExchangeItems() {
		if item.GetItemId() == itemId {
			return item
		}
	}
	return nil
}
