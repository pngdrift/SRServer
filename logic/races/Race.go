package races

import (
	"math/rand"
	"srserver/logic/db"
	"srserver/logic/playermanager"
	"srserver/logic/utils"
	pb "srserver/proto/out"
	"strconv"
)

func CreateEnemy(enemyPlayer *playermanager.Player) *pb.Enemy {
	if enemyPlayer == nil {
		return CreateBotEnemy()
	}
	return &pb.Enemy{
		Id:   enemyPlayer.User.Id,
		Info: enemyPlayer.User.Info,
		Car:  enemyPlayer.GetCurrentCar(), // or RandomCar?
		EnemyType: &pb.EnemyType{
			RaceCount:  0,
			Races:      []pb.RaceResult{},
			PlaceInTop: 0,
			Loot:       []*pb.CarUpgrade{},
		},
		Behavior: nil, // Automat
	}
}

var bots = map[int64]*pb.Enemy{}

func CreateBotEnemy() *pb.Enemy {
	botId := utils.MakeId()
	enemy := &pb.Enemy{
		Id: botId,
		Info: &pb.UserInfo{
			Name:   "Bot-" + strconv.FormatInt(botId, 16),
			Avatar: "https://robohash.org/" + strconv.FormatInt(botId, 16),
			Id:     botId,
			Lang:   "en",
			Type:   pb.UserType_TESTER,
		},
		Car: utils.BuildUserCar(db.Database.CarDatabase.Items[rand.Intn(len(db.Database.CarDatabase.Items))]),
		EnemyType: &pb.EnemyType{
			RaceCount:  0,
			Races:      []pb.RaceResult{},
			PlaceInTop: 0,
			Loot:       []*pb.CarUpgrade{},
		},
		Behavior: nil, // Automat
	}
	bots[botId] = enemy
	return enemy
}

func GetEnemyFor(player *playermanager.Player) *pb.Enemy {
	return CreateEnemy(playermanager.RandPlayer())
}

func GetEnemyById(enemyId int64) *pb.Enemy {
	enemyPlayer := playermanager.FindPlayerById(enemyId)
	if enemyPlayer != nil {
		return CreateEnemy(enemyPlayer)
	}
	return bots[enemyId]
}
