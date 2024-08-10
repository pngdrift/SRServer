//todo categorize

package utils

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	pb "srserver/proto/out"
	"time"
)

func BuildUserCar(baseCar *pb.BaseCar) *pb.UserCar {
	userCar := &pb.UserCar{}
	userCar.BaseId = baseCar.BaseId
	userCar.Id = MakeId()
	userCar.Number = &pb.CarNumber{
		CarId:     userCar.Id,
		IsTransit: true,
	}
	userCar.FrontBumperSlot = createUpgradeSlot(userCar.Id, baseCar.FrontBumperBaseId, pb.UpgradeType_FRONT_BUMPER)
	userCar.RearBumperSlot = createUpgradeSlot(userCar.Id, baseCar.RearBumperBaseId, pb.UpgradeType_REAR_BUMPER)
	userCar.CenterBumperSlot = createUpgradeSlot(userCar.Id, baseCar.CenterBumperBaseId, pb.UpgradeType_CENTER_BUMPER)

	// Maybe this is not needed
	// userCar.FrontSpringSlot = createUpgradeSlot(userCar.Id, baseCar.FrontSpringBaseId, pb.UpgradeType_SPRING)
	// userCar.RearSpringSlot = createUpgradeSlot(userCar.Id, baseCar.RearSpringBaseId, pb.UpgradeType_SPRING)
	// userCar.FrontBrakeSlot = createUpgradeSlot(userCar.Id, baseCar.FrontBrakeBaseId, pb.UpgradeType_BRAKE)
	// userCar.RearBrakeSlot = createUpgradeSlot(userCar.Id, baseCar.RearBrakeBaseId, pb.UpgradeType_BRAKE)
	// userCar.FrontSuspensionSlot = createUpgradeSlot(userCar.Id, baseCar.FrontSuspensionBaseId, pb.UpgradeType_FRONT_SUSPENSION_SUPPORT)
	// userCar.RearSuspensionSlot = createUpgradeSlot(userCar.Id, baseCar.RearSuspensionBaseId, pb.UpgradeType_REAR_SUSPENSION_SUPPORT)
	// userCar.DiskSlot = createUpgradeSlot(userCar.Id, baseCar.DiskBaseId, pb.UpgradeType_DISK)
	// userCar.FrontDiskSlot = createUpgradeSlot(userCar.Id, baseCar.DiskBaseId, pb.UpgradeType_DISK)
	// userCar.TiresSlot = createUpgradeSlot(userCar.Id, baseCar.TiresBaseId, pb.UpgradeType_TIRES)
	// userCar.FrontTiresSlot = createUpgradeSlot(userCar.Id, baseCar.TiresBaseId, pb.UpgradeType_TIRES)
	// userCar.TransmissionSlot = createUpgradeSlot(userCar.Id, baseCar.TransmissionBaseId, pb.UpgradeType_TRANSMISSION)
	userCar.EngineUpgrades = []*pb.EngineUpgrade{}
	for _, value := range pb.EngineUpgradeType_value {
		userCar.EngineUpgrades = append(userCar.EngineUpgrades, &pb.EngineUpgrade{
			Type:  pb.EngineUpgradeType(value),
			Level: 0,
		})
	}
	userCar.Paint = &pb.Paint{
		CarColor:           0,
		FrontBumperColor:   0,
		CenterBumperColor:  0,
		RearBumperColor:    0,
		DecalCounter:       0,
		Decals:             []*pb.Decal{},
		DiskColor:          0,
		IsDiskPainted:      false,
		TintingColor:       0,
		RimColor:           0,
		IsRimPainted:       false,
		Buyed:              false,
		DiskColorFront:     0,
		IsDiskPaintedFront: false,
		RimColorFront:      0,
		IsRimPaintedFront:  false,
	}
	userCar.Settings = &pb.CarSettings{
		Settings:       []*pb.Setting{},
		GearSettings:   []*pb.GearSetting{},
		MufflerOffsetX: 0,
		MufflerOffsetY: 0,
		YellowZoneRpm:  0,
		GreenZoneRpm:   0,
		RedZoneRpm:     0,
	}
	return userCar
}

func createUpgradeSlot(carId int64, baseId int32, upgradeType pb.UpgradeType) *pb.UpgradeSlot {
	return &pb.UpgradeSlot{
		CarId: carId,
		Type:  upgradeType,
		Upgrade: &pb.CarUpgrade{
			BaseId:   baseId,
			CarId:    carId,
			Current:  1.0,
			Grade:    pb.UpgradeGrade_WHITE,
			Id:       MakeId(),
			IsPacked: false,
			Type:     upgradeType,
		},
	}
}

func GetDefaultUser() *pb.User {
	user := &pb.User{}

	user.Level = 50
	user.Exp = 0
	user.Money = &pb.Money{
		Gold:             10000000,
		Money:            10000000,
		TopPoints:        10000000,
		TournamentPoints: 10000000,
		UpgradePoints:    10000000,
	}
	user.Fuel = &pb.Fuel{
		Addition: 0,
		Fuel:     75,
		FuelTime: int32(time.Now().Unix()),
	}
	user.Info = &pb.UserInfo{
		Avatar: "",
		Id:     user.Id,
		Lang:   "en",
		Name:   "Aivaz",
		Type:   pb.UserType_ADMIN,
	}
	user.Garage = &pb.Garage{
		CurrentId: 0,
		Cars:      []*pb.UserCar{},
	}
	user.Inventory = &pb.Inventory{
		Items:    []*pb.Item{},
		Upgrades: []*pb.CarUpgrade{},
	}
	user.Enemies = &pb.UserEnemies{
		CarId: 0,
		List:  []*pb.Enemy{},
	}
	user.Quests = &pb.UserQuests{
		Quests: []*pb.Quest{
			{
				BaseId:     7,
				Counter:    0,
				Desc:       "Buy lol",
				IsFinished: false,
				Name:       "Railgun master",
				SaveTime:   10000,
			},
		},
	}
	user.World = &pb.World{
		Time:      0,
		TimeDelta: 0,
		ZoneId:    "default",
	}
	user.Mail = &pb.MailBox{
		IsLoaded: true,
		Mails: []*pb.MailMessage{
			{
				Id:       MakeId(),
				FromName: "System",
				FromUid:  0,
				ToUid:    0,
				Time:     time.Now().UnixMilli(),
				Title:    "Welcome to Street racing!",
				Message:  "The first of its kind racing game built on a physics engine with a wide customization.",
				IsReaded: false,
				IsSystem: true,
				Money: &pb.Money{
					Gold:  200,
					Money: 10000,
				},
				Exp:      0,
				Fuel:     0,
				Upgrades: []*pb.CarUpgrade{},
				Items:    []*pb.Item{},
			},
		},
	}
	user.DailyBonusTime = 0
	user.PointsEnemies = &pb.PointsEnemies{
		IsNeedUpdate: true,
		List:         []*pb.Enemy{},
	}
	user.TimersAndCounters = &pb.TimersAndCounters{
		ResetTime: time.Now().UnixMilli(),
	}
	user.Challenges = &pb.Challenges{
		Items: []*pb.ChallengeItem{},
	}
	user.SocialType = pb.SocialType_DEBUG
	user.Dyno = &pb.Dyno{}
	user.Paints = &pb.UserPaints{
		Items: []*pb.PaintItem{},
	}

	return user
}

func GetBotEnemy() *pb.Enemy {
	return &pb.Enemy{
		Id: MakeId(),
		Info: &pb.UserInfo{
			Avatar: "https://v-item.ru/assets/cars/car_11.pnG",
			Id:     MakeId(),
			Lang:   "RU",
			Name:   "Bot",
			Type:   pb.UserType_TESTER,
		},
		Car: &pb.UserCar{
			Id:     1,
			BaseId: 3,
			Number: &pb.CarNumber{
				IsTransit: true,
			},
		},
		Type: pb.RaceType_RACE,
		EnemyType: &pb.EnemyType{
			RaceCount:  0,
			Races:      []pb.RaceResult{},
			PlaceInTop: 0,
			Loot:       []*pb.CarUpgrade{},
		},
		Behavior: &pb.Behavior{
			TiresHeat:        0.7,
			TransmissionType: 1,
			FrontTiresHeat:   1,
			RearTiresHeat:    1,
			StartRpm:         2000,
			Events: []*pb.RaceEventItem{
				{
					Time:  0,
					Event: pb.RaceEvent_GAS_DOWN,
				},
				{
					Time:  2,
					Event: pb.RaceEvent_SHIFT_UP,
				},
				{
					Time:  3,
					Event: pb.RaceEvent_SHIFT_UP,
				},
				{
					Time:  6,
					Event: pb.RaceEvent_SHIFT_UP,
				},
			},
		},
	}
}

func MakeId() int64 {
	// Even though id is int64, but if you generate a number
	// greater than int32, then there are may some bugs in game
	return int64(rand.Int31())
}

func LoadAvatar(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error while making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body: %v", err)
	}

	return imageBytes, nil
}

func MultiplyMoney(money *pb.Money, count int32) *pb.Money {
	money.Gold *= count
	money.Money *= count
	money.TopPoints *= count
	money.TournamentPoints *= count
	money.UpgradePoints *= count
	return money
}

func GetUpgradeType(slotName string) pb.UpgradeType {
	switch slotName {
	case "FRONT_DISK_SLOT", "DISK_SLOT":
		return pb.UpgradeType_DISK
	case "FRONT_BRAKE_SLOT", "REAR_BRAKE_SLOT":
		return pb.UpgradeType_BRAKE
	case "FRONT_SPRING_SLOT", "REAR_SPRING_SLOT":
		return pb.UpgradeType_SPRING
	case "FRONT_SUSPENSION_SLOT", "REAR_SUSPENSION_SLOT":
		return pb.UpgradeType_SUSPENSION
	case "TIRES_SLOT", "FRONT_TIRES_SLOT":
		return pb.UpgradeType_TIRES
	case "CENTER_BUMPER_SLOT":
		return pb.UpgradeType_CENTER_BUMPER
	case "FRONT_WHEEL_SLOT", "WHEEL_SLOT":
		return pb.UpgradeType_WHEEL_PART
	case "FRAME_SLOT":
		return pb.UpgradeType_FRAME_PART
	case "REAR_BRAKE_PAD_SLOT", "FRONT_BRAKE_PAD_SLOT":
		return pb.UpgradeType_BRAKE_PAD
	case "FRONT_BUMPER_SLOT":
		return pb.UpgradeType_FRONT_BUMPER
	case "HOOD_SLOT":
		return pb.UpgradeType_HOOD_PART
	case "PNEUMO_SLOT":
		return pb.UpgradeType_PNEUMATIC_SUSPENSION
	case "REAR_BUMPER_SLOT":
		return pb.UpgradeType_REAR_BUMPER
	case "ROOF_SLOT":
		return pb.UpgradeType_ROOF_PART
	case "TRUNK_SLOT":
		return pb.UpgradeType_TRUNK_PART
	case "TURBO_1_SLOT":
		return pb.UpgradeType_TURBO_1
	case "TURBO_2_SLOT":
		return pb.UpgradeType_TURBO_2
	case "HEADLIGHT_SLOT":
		return pb.UpgradeType_HEADLIGHT
	case "NEON_DISK_SLOT":
		return pb.UpgradeType_NEON_DISK
	case "NEON_SLOT":
		return pb.UpgradeType_NEON
	case "SPOILER_SLOT":
		return pb.UpgradeType_SPOILER
	case "TRANSMISSION_SLOT":
		return pb.UpgradeType_TRANSMISSION
	case "DIFFERENTIAL_SLOT":
		return pb.UpgradeType_DIFFERENTIAL
	case "AIR_FILTER_SLOT":
		return pb.UpgradeType_AIR_FILTER
	case "INTERCOOLER_SLOT":
		return pb.UpgradeType_INTERCOOLER
	case "PIPE_SLOT":
		return pb.UpgradeType_PIPES
	case "INTAKE_MAINFOLD_SLOT":
		return pb.UpgradeType_INTAKE_MAINFOLD
	case "WESTGATE_SLOT":
		return pb.UpgradeType_WESTGATE
	case "EXHAUST_MAINFOLD_SLOT":
		return pb.UpgradeType_EXHAUST_MAINFOLD
	case "EXHAUST_OUTLET_SLOT":
		return pb.UpgradeType_EXHAUST_OUTLET
	case "EXHAUST_MUFFLER_SLOT":
		return pb.UpgradeType_EXHAUST_MUFFLER
	case "TIMING_GEAR_SLOT":
		return pb.UpgradeType_TIMING_GEAR
	case "CAMSHAFT_SLOT":
		return pb.UpgradeType_CAMSHAFTS
	case "ECU_SLOT":
		return pb.UpgradeType_ECU
	default:
		return pb.UpgradeType_DUMMY
	}
}

func InstallUpgrade(userCar *pb.UserCar, upgradeSlot *pb.UpgradeSlot, slotName string) {
	switch slotName {
	case "FRONT_DISK_SLOT":
		userCar.FrontDiskSlot = upgradeSlot
	case "DISK_SLOT":
		userCar.DiskSlot = upgradeSlot
	case "FRONT_BRAKE_SLOT":
		userCar.FrontBrakeSlot = upgradeSlot
	case "REAR_BRAKE_SLOT":
		userCar.RearBrakeSlot = upgradeSlot
	case "FRONT_SPRING_SLOT":
		userCar.FrontSpringSlot = upgradeSlot
	case "REAR_SPRING_SLOT":
		userCar.RearSpringSlot = upgradeSlot
	case "FRONT_SUSPENSION_SLOT":
		userCar.FrontSuspensionSlot = upgradeSlot
	case "REAR_SUSPENSION_SLOT":
		userCar.RearSuspensionSlot = upgradeSlot
	case "TIRES_SLOT":
		userCar.TiresSlot = upgradeSlot
	case "FRONT_TIRES_SLOT":
		userCar.FrontTiresSlot = upgradeSlot
	case "CENTER_BUMPER_SLOT":
		userCar.CenterBumperSlot = upgradeSlot
	case "FRONT_WHEEL_SLOT":
		userCar.FrontWheelSlot = upgradeSlot
	case "WHEEL_SLOT":
		userCar.WheelSlot = upgradeSlot
	case "FRAME_SLOT":
		userCar.FrameSlot = upgradeSlot
	case "REAR_BRAKE_PAD_SLOT":
		userCar.RearBrakePadSlot = upgradeSlot
	case "FRONT_BRAKE_PAD_SLOT":
		userCar.FrontBrakePadSlot = upgradeSlot
	case "FRONT_BUMPER_SLOT":
		userCar.FrontBumperSlot = upgradeSlot
	case "HOOD_SLOT":
		userCar.HoodSlot = upgradeSlot
	case "PNEUMO_SLOT":
		userCar.PneumoSlot = upgradeSlot
	case "REAR_BUMPER_SLOT":
		userCar.RearBumperSlot = upgradeSlot
	case "ROOF_SLOT":
		userCar.RoofSlot = upgradeSlot
	case "TRUNK_SLOT":
		userCar.TrunkSlot = upgradeSlot
	case "TURBO_1_SLOT":
		userCar.Turbo1Slot = upgradeSlot
	case "TURBO_2_SLOT":
		userCar.Turbo2Slot = upgradeSlot
	case "HEADLIGHT_SLOT":
		userCar.HeadlightSlot = upgradeSlot
	case "NEON_DISK_SLOT":
		userCar.NeonDiskSlot = upgradeSlot
	case "NEON_SLOT":
		userCar.NeonSlot = upgradeSlot
	case "SPOILER_SLOT":
		userCar.SpoilerSlot = upgradeSlot
	case "TRANSMISSION_SLOT":
		userCar.TransmissionSlot = upgradeSlot
	case "DIFFERENTIAL_SLOT":
		userCar.DifferentialSlot = upgradeSlot
	case "AIR_FILTER_SLOT":
		userCar.AirFilterSlot = upgradeSlot
	case "INTERCOOLER_SLOT":
		userCar.IntercoolerSlot = upgradeSlot
	case "PIPE_SLOT":
		userCar.PipeSlot = upgradeSlot
	case "INTAKE_MAINFOLD_SLOT":
		userCar.IntakeMainfoldSlot = upgradeSlot
	case "WESTGATE_SLOT":
		userCar.WestgateSlot = upgradeSlot
	case "EXHAUST_MAINFOLD_SLOT":
		userCar.ExhaustMainfoldSlot = upgradeSlot
	case "EXHAUST_OUTLET_SLOT":
		userCar.ExhaustOutletSlot = upgradeSlot
	case "EXHAUST_MUFFLER_SLOT":
		userCar.ExhaustMufflerSlot = upgradeSlot
	case "TIMING_GEAR_SLOT":
		userCar.TimingGearSlot = upgradeSlot
	case "CAMSHAFT_SLOT":
		userCar.CamshaftSlot = upgradeSlot
	case "ECU_SLOT":
		userCar.EcuSlot = upgradeSlot
	}
}

func GetUpgradeSlot(userCar *pb.UserCar, slotName string) *pb.UpgradeSlot {
	switch slotName {
	case "FRONT_DISK_SLOT":
		return userCar.FrontDiskSlot
	case "DISK_SLOT":
		return userCar.DiskSlot
	case "FRONT_BRAKE_SLOT":
		return userCar.FrontBrakeSlot
	case "REAR_BRAKE_SLOT":
		return userCar.RearBrakeSlot
	case "FRONT_SPRING_SLOT":
		return userCar.FrontSpringSlot
	case "REAR_SPRING_SLOT":
		return userCar.RearSpringSlot
	case "FRONT_SUSPENSION_SLOT":
		return userCar.FrontSuspensionSlot
	case "REAR_SUSPENSION_SLOT":
		return userCar.RearSuspensionSlot
	case "TIRES_SLOT":
		return userCar.TiresSlot
	case "FRONT_TIRES_SLOT":
		return userCar.FrontTiresSlot
	case "CENTER_BUMPER_SLOT":
		return userCar.CenterBumperSlot
	case "FRONT_WHEEL_SLOT":
		return userCar.FrontWheelSlot
	case "WHEEL_SLOT":
		return userCar.WheelSlot
	case "FRAME_SLOT":
		return userCar.FrameSlot
	case "REAR_BRAKE_PAD_SLOT":
		return userCar.RearBrakePadSlot
	case "FRONT_BRAKE_PAD_SLOT":
		return userCar.FrontBrakePadSlot
	case "FRONT_BUMPER_SLOT":
		return userCar.FrontBumperSlot
	case "HOOD_SLOT":
		return userCar.HoodSlot
	case "PNEUMO_SLOT":
		return userCar.PneumoSlot
	case "REAR_BUMPER_SLOT":
		return userCar.RearBumperSlot
	case "ROOF_SLOT":
		return userCar.RoofSlot
	case "TRUNK_SLOT":
		return userCar.TrunkSlot
	case "TURBO_1_SLOT":
		return userCar.Turbo1Slot
	case "TURBO_2_SLOT":
		return userCar.Turbo2Slot
	case "HEADLIGHT_SLOT":
		return userCar.HeadlightSlot
	case "NEON_DISK_SLOT":
		return userCar.NeonDiskSlot
	case "NEON_SLOT":
		return userCar.NeonSlot
	case "SPOILER_SLOT":
		return userCar.SpoilerSlot
	case "TRANSMISSION_SLOT":
		return userCar.TransmissionSlot
	case "DIFFERENTIAL_SLOT":
		return userCar.DifferentialSlot
	case "AIR_FILTER_SLOT":
		return userCar.AirFilterSlot
	case "INTERCOOLER_SLOT":
		return userCar.IntercoolerSlot
	case "PIPE_SLOT":
		return userCar.PipeSlot
	case "INTAKE_MAINFOLD_SLOT":
		return userCar.IntakeMainfoldSlot
	case "WESTGATE_SLOT":
		return userCar.WestgateSlot
	case "EXHAUST_MAINFOLD_SLOT":
		return userCar.ExhaustMainfoldSlot
	case "EXHAUST_OUTLET_SLOT":
		return userCar.ExhaustOutletSlot
	case "EXHAUST_MUFFLER_SLOT":
		return userCar.ExhaustMufflerSlot
	case "TIMING_GEAR_SLOT":
		return userCar.TimingGearSlot
	case "CAMSHAFT_SLOT":
		return userCar.CamshaftSlot
	case "ECU_SLOT":
		return userCar.EcuSlot
	default:
		return nil
	}
}

func BuildTrack(baseTrack *pb.BaseTrack, zoneId string) *pb.Track {
	track := &pb.Track{}
	track.BaseId = baseTrack.BaseId
	track.BrakeDistance = 200
	track.GroundWidth = baseTrack.GroundWidth
	track.GroundHeight = baseTrack.GroundHeight
	track.GroundOffset = baseTrack.GroundOffset
	track.GroundFriction = baseTrack.GroundFriction
	track.GroundStep = baseTrack.GroundStep
	track.Frequency = baseTrack.Frequency
	track.OctaveCount = baseTrack.OctaveCount
	track.Lacunarity = baseTrack.Lacunarity
	track.Persistence = baseTrack.Persistence
	track.FinishLineY1 = baseTrack.FinishLineY1
	track.FinishLineY2 = baseTrack.FinishLineY2

	if baseTrack.AutoTimesOfDay {
		dayPostfix := "_" + GetCurrentDayState(zoneId)
		track.Ground = baseTrack.GetGround() + dayPostfix
		track.Backgrounds = make([]*pb.BaseTrackBackground, len(baseTrack.Backgrounds))
		for i, background := range baseTrack.Backgrounds {
			track.Backgrounds[i] = &pb.BaseTrackBackground{
				Name:         background.Name + dayPostfix,
				OffsetFactor: background.OffsetFactor,
				Width:        background.Width,
				Height:       background.Height,
				Offset:       background.Offset,
			}
		}
	} else {
		track.Backgrounds = baseTrack.Backgrounds
	}
	return track
}

func GetCurrentDayState(zoneId string) string {
	loc, err := time.LoadLocation(zoneId)
	if err != nil {
		loc = time.UTC
	}
	currentHour := time.Now().In(loc).Hour()
	switch {
	case currentHour >= 6 && currentHour < 12:
		return "morning"
	case currentHour >= 12 && currentHour < 18:
		return "day"
	case currentHour >= 18 && currentHour < 21:
		return "twilight"
	default:
		return "night"
	}
}
