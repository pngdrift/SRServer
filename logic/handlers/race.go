package handlers

import (
	"log"
	"math/rand"
	"srserver/common"
	"srserver/conf"
	"srserver/logic/clients"
	"srserver/logic/db"
	"srserver/logic/races"
	"srserver/logic/utils"
	pb "srserver/proto/out"
	"time"
)

// todo
func GetRaceAward(pack *common.Pack, client *clients.Client) {
	finishParams := &pb.FinishParams{}
	pack.ReadProto(finishParams)
	player := client.Player
	//player.GetCurrentCar().Behavior = finishParams.Behavior
	player.GetCurrentCar().AccumulatorDistance += finishParams.UserDistance
	response := common.NewResponse(pack)
	race := &pb.Race{}
	race.Money = &pb.Money{
		Money: 1000,
	}
	race.Items = []*pb.Item{
		{
			Id:     utils.MakeId(),
			Count:  1,
			Type:   pb.ItemType_BLUEPRINT,
			BaseId: 101,
		},
	}
	var result pb.RaceResult
	if finishParams.UserTime < finishParams.EnemyTime {
		result = pb.RaceResult_WIN
	} else if (finishParams.UserTime - finishParams.EnemyTime) < 0.001 {
		result = pb.RaceResult_DRAW
	} else {
		result = pb.RaceResult_LOST
	}

	race.Result = result

	if finishParams.EnemyId > 0 {
		enemy := races.GetEnemyById(finishParams.EnemyId)
		enemy.EnemyType.RaceCount++
		enemy.EnemyType.Races = append(enemy.EnemyType.Races, result)

		if enemy.EnemyType.RaceCount < 3 {
			race.IsCanRepeat = true
		}
	}

	race.Type = finishParams.Type
	race.UserTime = finishParams.UserTime

	response.WriteProto(race)
	client.Write(response.ToByteBuff())
}

func UpdateEnemies(pack *common.Pack, client *clients.Client) {
	enemies := &pb.UserEnemies{
		CarId: client.Player.User.Garage.CurrentId,
		List: []*pb.Enemy{
			races.CreateBotEnemy(),
			races.CreateBotEnemy(),
			//race.GetEnemyFor(client.Player),
			races.CreateBotEnemy(),
		},
	}
	client.Player.User.Enemies = enemies
	client.Player.Save()
	response := common.NewResponse(pack)
	response.WriteProto(enemies)
	client.Write(response.ToByteBuff())
}

func CreateRace(pack *common.Pack, client *clients.Client) {
	startParams := &pb.StartParams{}
	pack.ReadProto(startParams)
	player := client.Player
	tc := player.User.TimersAndCounters
	var enemy *pb.Enemy
	var distance float32
	var baseTrack *pb.BaseTrack = db.FindTrack(2018)
	switch startParams.Type {
	case pb.RaceType_TEST402:
		distance = 402
		player.Withdraw(conf.MONEY_FOR_TEST_RACE)
	case pb.RaceType_TEST804:
		distance = 804
		player.Withdraw(conf.MONEY_FOR_TEST_RACE)
	case pb.RaceType_SHADOW:
		distance = 804
		player.Withdraw(conf.MONEY_FOR_TEST_RACE)
		enemy = races.CreateEnemy(player)
	case pb.RaceType_RACE:
		if !player.IsFullEnought(5) {
			return
		}
		player.WithdrawFuel(5)
		distance = 402
		baseTrack = nil
		enemy = races.GetEnemyById(startParams.EnemyId)
	case pb.RaceType_POINTS:
		if tc.RatingCount >= 5 || tc.RatingTimer > time.Now().UnixMilli() {
			return
		}
		tc.RatingCount++
		tc.RatingTimer = time.Now().UnixMilli() + conf.RATING_RACE_DELAY
		distance = 804
		enemy = races.GetEnemyById(startParams.EnemyId)
	case pb.RaceType_TIME:
		if tc.TimeCount >= 5 || tc.TimeTimer > time.Now().UnixMilli() {
			return
		}
		tc.TimeCount++
		tc.TimeTimer = time.Now().UnixMilli() + conf.TIME_RACE_DELAY
		distance = 804
	case pb.RaceType_CHALLENGE:
		if tc.ChallengeCount >= 5 || tc.ChallengeTimer > time.Now().UnixMilli() {
			return
		}
		tc.ChallengeCount++
		tc.ChallengeTimer = time.Now().UnixMilli() + conf.CHALLENGE_RACE_DELAY
		distance = 402

	default:
		log.Print("Unknown race type ", startParams.Type)
		return
	}

	if baseTrack == nil {
		baseTracks := db.Database.TrackDatabase.GetItems()
		baseTrack = baseTracks[rand.Intn(len(baseTracks))]
	}
	track := utils.BuildTrack(baseTrack, player.User.World.ZoneId)
	track.Distanse = distance
	track.Distanse += 3.25 // Burnout zone distance

	response := common.NewResponse(pack)
	response.WriteProto(track)
	if enemy != nil {
		enemy.Type = startParams.Type
		response.WriteProto(enemy)
	}
	client.Write(response.ToByteBuff())

}

func BrakeRace(pack *common.Pack, client *clients.Client) {
	// idk why for BrakeRace take fuel (client hardcoded)
	client.Player.WithdrawFuel(5)
}

func UpdatePointsEnemies(pack *common.Pack, client *clients.Client) {
	pointsEnemies := &pb.PointsEnemies{
		List: []*pb.Enemy{
			races.CreateBotEnemy(),
			races.CreateBotEnemy(),
		},
		IsNeedUpdate: false,
	}
	client.Player.User.PointsEnemies = pointsEnemies
	response := common.NewResponse(pack)
	response.WriteProto(pointsEnemies)
	client.Write(response.ToByteBuff())
}
