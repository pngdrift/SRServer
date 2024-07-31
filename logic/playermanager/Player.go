package playermanager

import (
	"log"
	"math/rand"
	"net"
	"srserver/common"
	"srserver/conf"
	"srserver/logic/utils"
	pb "srserver/proto/out"
	"time"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	Id   int64
	User *pb.User
	Uid  string
	Conn net.Conn
}

func (p *Player) Save() error {
	query := `
		INSERT INTO players (uid, userProto) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE userProto=VALUES(userProto)`
	userData, err := proto.Marshal(p.User)
	if err != nil {
		return err
	}
	_, err = GetDb().Exec(query, p.Uid, userData)
	return err
}

func (p *Player) BuyCar(car *pb.BaseCar) *pb.UserCar {
	userCar := utils.BuildUserCar(car)
	p.User.GetGarage().Cars = append(p.User.GetGarage().Cars, userCar)
	p.User.GetGarage().CurrentId = userCar.Id
	p.Withdraw(car.GetPrice())
	return userCar
}

func (p *Player) GetCurrentCar() *pb.UserCar {
	return p.FindUserCarById(p.User.Garage.CurrentId)
}

func (p *Player) GetRandomCar() *pb.UserCar {
	return p.User.Garage.Cars[rand.Intn(len(p.User.Garage.Cars))]
}

func (p *Player) Withdraw(money *pb.Money) {
	p.User.Money.Money -= money.Money
	p.User.Money.TopPoints -= money.TopPoints
	p.User.Money.Gold -= money.Gold
	p.User.Money.UpgradePoints -= money.UpgradePoints
	p.User.Money.TournamentPoints -= money.TournamentPoints
	p.Save()
}

func (p *Player) Deposit(money *pb.Money) {
	p.User.Money.Money += money.Money
	p.User.Money.TopPoints += money.TopPoints
	p.User.Money.Gold += money.Gold
	p.User.Money.UpgradePoints += money.UpgradePoints
	p.User.Money.TournamentPoints += money.TournamentPoints
	p.Save()
}

func (p *Player) CheckMoney(money *pb.Money) bool {
	return p.User.Money.Money >= money.Money &&
		p.User.Money.TopPoints >= money.TopPoints &&
		p.User.Money.Gold >= money.Gold &&
		p.User.Money.UpgradePoints >= money.UpgradePoints &&
		p.User.Money.TournamentPoints >= money.TournamentPoints
}

func (p *Player) getExpToLevel() int32 {
	return (p.User.Level * 10) + 50
}

func (p *Player) AddExp(exp int32) int32 {
	if exp <= 0 {
		return 0
	}
	expToLevel := p.getExpToLevel()
	if p.User.Exp+exp >= expToLevel {
		log.Println("Level up!")
		p.User.Level++

		award := &pb.LevelUpAward{
			Level: p.User.Level,
			Money: &pb.Money{
				Gold: 200,
			},
			Fuel:     &pb.Fuel{},
			Upgrades: []*pb.CarUpgrade{},
		}
		p.Deposit(award.Money)

		pack := common.NewPack()
		pack.SetMethod(common.OnLevelUp)
		pack.SetSequence(10)
		pack.WriteProto(award)
		p.Conn.Write(pack.ToByteBuff())

		deltaExp := expToLevel - p.User.Exp
		p.User.Exp = 0
		levelUp := p.AddExp(exp-deltaExp) + 1
		return levelUp
	}
	p.User.Exp += exp
	p.Save()
	return 0
}

func (p *Player) FindUserCarById(id int64) *pb.UserCar {
	for _, car := range p.User.GetGarage().GetCars() {
		if car.Id == id {
			return car
		}
	}
	return nil
}

func (p *Player) FindPaintItemById(id int64) *pb.PaintItem {
	for _, item := range p.User.Paints.Items {
		if item.Id == id {
			return item
		}
	}
	return nil
}

func (p *Player) updateFuel() {
	currentTime := int32(time.Now().Unix())
	deltaTime := currentTime - p.User.Fuel.FuelTime
	fuelToRestore := deltaTime / conf.FUEL_RESTORE_TIME
	elapsedTime := deltaTime % conf.FUEL_RESTORE_TIME
	if fuelToRestore > 0 {
		p.User.Fuel.Fuel += fuelToRestore
		p.User.Fuel.FuelTime = currentTime + elapsedTime
	}
	if p.User.Fuel.Fuel > conf.MAX_FUEL {
		p.User.Fuel.Fuel = conf.MAX_FUEL
		p.User.Fuel.FuelTime = currentTime
	}
}

func (p *Player) AddFuel(value int32) {
	p.updateFuel()
	fuel := p.User.Fuel
	fuel.Fuel += value
	if fuel.Fuel > conf.MAX_FUEL {
		fuel.Addition += fuel.Fuel - conf.MAX_FUEL
		fuel.Fuel = conf.MAX_FUEL
	}
	p.Save()
}

func (p *Player) WithdrawFuel(value int32) {
	p.updateFuel()
	fuel := p.User.Fuel
	totalFuel := fuel.Fuel + fuel.Addition
	remainTotalFuel := totalFuel - value
	if remainTotalFuel > conf.MAX_FUEL {
		fuel.Addition = remainTotalFuel - conf.MAX_FUEL
		fuel.Fuel = conf.MAX_FUEL
	} else {
		fuel.Addition = 0
		fuel.Fuel = remainTotalFuel
	}
	p.Save()
}

func (p *Player) IsFullEnought(value int32) bool {
	fuel := p.User.Fuel
	totalFuel := fuel.Fuel + fuel.Addition
	return totalFuel >= value
}

func (p *Player) CheckResetTime() {
	currentMillis := time.Now().UnixMilli()
	tc := p.User.TimersAndCounters
	if currentMillis >= tc.ResetTime {
		tc.TimeCount = 0
		tc.TimeTimer = 0
		tc.RatingCount = 0
		tc.RatingTimer = 0
		tc.ChallengeCount = 0
		tc.ChallengeTimer = 0
		tc.ExchangeCount = 0
		tc.ExchangeTimer = 0
		tc.RefuelCount = 0
		tc.ResetTime = currentMillis + 1000*60*60*18 // +18 hours
		p.Save()
	}
}
