package conf

import pb "srserver/proto/out"

// Vars that is hardcoded in the game
var (
	MAX_FUEL int32 = 75

	FUEL_RESTORE_TIME int32 = 180

	REFUEL_COST = []*pb.Money{{
		Gold: 20,
	}, {
		Gold: 50,
	}, {
		Gold: 100,
	}, {
		Gold: 100,
	}, {
		Gold: 200,
	}}

	REFUEL_VALUE int32 = 100

	UPGRADE_COST_MAP = map[pb.EngineUpgradeType][]*pb.Money{
		pb.EngineUpgradeType_GEARS:         {{Money: 0, Gold: 0}, {Money: 300, Gold: 0}, {Money: 850, Gold: 0}, {Money: 3200, Gold: 0}, {Money: 0, Gold: 40}, {Money: 0, Gold: 60}, {Money: 11900, Gold: 0}, {Money: 11900, Gold: 0}, {Money: 0, Gold: 90}, {Money: 0, Gold: 110}, {Money: 25200, Gold: 0}},
		pb.EngineUpgradeType_EXHAUST:       {{Money: 0, Gold: 0}, {Money: 850, Gold: 0}, {Money: 900, Gold: 0}, {Money: 3550, Gold: 0}, {Money: 0, Gold: 20}, {Money: 0, Gold: 50}, {Money: 15000, Gold: 0}, {Money: 15233, Gold: 0}, {Money: 0, Gold: 100}, {Money: 0, Gold: 110}, {Money: 21300, Gold: 0}},
		pb.EngineUpgradeType_CANDLE:        {{Money: 0, Gold: 0}, {Money: 450, Gold: 0}, {Money: 1790, Gold: 0}, {Money: 4100, Gold: 0}, {Money: 0, Gold: 30}, {Money: 0, Gold: 60}, {Money: 12000, Gold: 0}, {Money: 12600, Gold: 0}, {Money: 0, Gold: 110}, {Money: 0, Gold: 110}, {Money: 25750, Gold: 0}},
		pb.EngineUpgradeType_PISTON:        {{Money: 0, Gold: 0}, {Money: 900, Gold: 0}, {Money: 550, Gold: 0}, {Money: 1910, Gold: 0}, {Money: 0, Gold: 40}, {Money: 0, Gold: 70}, {Money: 15800, Gold: 0}, {Money: 15700, Gold: 0}, {Money: 0, Gold: 80}, {Money: 0, Gold: 130}, {Money: 29600, Gold: 0}},
		pb.EngineUpgradeType_ROD:           {{Money: 0, Gold: 0}, {Money: 350, Gold: 0}, {Money: 1090, Gold: 0}, {Money: 3150, Gold: 0}, {Money: 0, Gold: 50}, {Money: 0, Gold: 80}, {Money: 11000, Gold: 0}, {Money: 11620, Gold: 0}, {Money: 0, Gold: 130}, {Money: 0, Gold: 130}, {Money: 27650, Gold: 0}},
		pb.EngineUpgradeType_CYLINDER_HEAD: {{Money: 0, Gold: 0}, {Money: 120, Gold: 0}, {Money: 800, Gold: 0}, {Money: 2990, Gold: 0}, {Money: 0, Gold: 30}, {Money: 0, Gold: 40}, {Money: 15100, Gold: 0}, {Money: 15400, Gold: 0}, {Money: 0, Gold: 150}, {Money: 0, Gold: 180}, {Money: 22300, Gold: 0}},
		pb.EngineUpgradeType_CAMSHAFT:      {{Money: 0, Gold: 0}, {Money: 300, Gold: 0}, {Money: 2100, Gold: 0}, {Money: 3860, Gold: 0}, {Money: 0, Gold: 40}, {Money: 0, Gold: 80}, {Money: 15200, Gold: 0}, {Money: 15700, Gold: 0}, {Money: 0, Gold: 110}, {Money: 0, Gold: 140}, {Money: 22600, Gold: 0}},
		pb.EngineUpgradeType_FUEL_PUMP:     {{Money: 0, Gold: 0}, {Money: 450, Gold: 0}, {Money: 1800, Gold: 0}, {Money: 4250, Gold: 0}, {Money: 0, Gold: 20}, {Money: 0, Gold: 60}, {Money: 11000, Gold: 0}, {Money: 11650, Gold: 0}, {Money: 0, Gold: 100}, {Money: 0, Gold: 180}, {Money: 27800, Gold: 0}},
	}

	MONEY_FOR_TEST_RACE = &pb.Money{Money: 400}

	MONEY_FOR_DYNO_TEST = &pb.Money{Money: 300}

	SHOP_NUMBERS_COST = &pb.Money{Money: 5000, Gold: 50}

	TIME_RACE_DELAY int64 = 300 * 1000

	RATING_RACE_DELAY int64 = 300 * 1000

	CHALLENGE_RACE_DELAY int64 = 300 * 1000

	MIN_MUFFLER_OFFSET_X float32 = -0.05
	MAX_MUFFLER_OFFSET_X float32 = 0.10

	MIN_MUFFLER_OFFSET_Y float32 = -0.05
	MAX_MUFFLER_OFFSET_Y float32 = 0.15
)
