package db

import (
	"hash/crc32"
	"io"
	"log"
	"os"
	"path/filepath"
	pb "srserver/proto/out"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var Database *pb.Database
var DatabaseBuff []byte
var DatabaseChecksum int64

func InitDb() {
	// Combining jsons into one
	var json = "{"
	files, err := filepath.Glob("./resources/db/*Database.json")
	if err != nil {
		log.Fatalf("Error finding files: %v", err)
	}
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("Error opening file %s: %v", file, err)
		}
		content, _ := io.ReadAll(f)
		f.Close()
		contentStr := string(content)
		offset := 1
		if contentStr[len(contentStr)-1] != '}' {
			offset = 2
		}
		contentStr = contentStr[1 : len(contentStr)-offset]
		json += contentStr + ","
	}

	req := &pb.Database{}
	err = protojson.Unmarshal([]byte(json[0:len(json)-1]+"}"), req)
	if err != nil {
		log.Fatalln("Failed Unmarshal:", err)
	}

	carsMap := map[int32]*pb.BaseCar{}
	for _, baseCar := range req.CarDatabase.Items {
		carsMap[baseCar.BaseId] = baseCar
		_, err := os.Stat("./resources/assets-ext/" + baseCar.Image)
		if err != nil {
			log.Fatal("Car atlas not found: ", baseCar.Image)
		}
	}

	for _, tuning := range req.TuningDatabase.Items {
		if _, exists := carsMap[tuning.BaseCarId]; !exists && tuning.BaseCarId > 0 {
			log.Fatal("Found tuning for unexisting car:", tuning.BaseCarId)
		}
	}

	req.ChallengeDatabase = createChallengeDatabase()
	req.BlueprintDatabase = createBlueprintsDatabase()
	Database = req
	out, err := proto.Marshal(Database)
	if err != nil {
		log.Fatalln("Failed to encode database:", err)
	}
	DatabaseBuff = out
	DatabaseChecksum = int64(crc32.ChecksumIEEE(out))
	log.Printf("Objects database created: checksum: %s size: %vKB", strconv.FormatInt(DatabaseChecksum, 16), len(DatabaseBuff)/1024)
}

func createChallengeDatabase() *pb.ChallengeDatabase {
	result := &pb.ChallengeDatabase{
		Items: []*pb.BaseChallenge{},
	}
	for i := 1; i <= 7; i++ {
		result.Items = append(result.Items, &pb.BaseChallenge{
			BaseId:    int32(i),
			Day:       int32(i),
			Classes:   []string{"A", "B", "C", "D", "E", "F", "G", "I", "J", "K", "L"},
			DriveType: pb.DriveType_FRONT,
			Image:     "challenge_id" + strconv.Itoa(i) + "_icon",
		})
	}
	return result
}

func createBlueprintsDatabase() *pb.BlueprintDatabase {
	result := &pb.BlueprintDatabase{
		Items: []*pb.BaseBlueprint{},
	}
	upgradeTypes := map[pb.UpgradeType]bool{
		pb.UpgradeType_HOOD_PART:  true,
		pb.UpgradeType_TRUNK_PART: true,
		pb.UpgradeType_ROOF_PART:  true,
		pb.UpgradeType_FRAME_PART: true,

		pb.UpgradeType_DISK:          false,
		pb.UpgradeType_TIRES:         false,
		pb.UpgradeType_BRAKE:         false,
		pb.UpgradeType_BRAKE_PAD:     false,
		pb.UpgradeType_FRONT_BUMPER:  false,
		pb.UpgradeType_REAR_BUMPER:   false,
		pb.UpgradeType_CENTER_BUMPER: false,
		pb.UpgradeType_SPOILER:       false,
		//pb.UpgradeType_HEADLIGHT:        false,
		//pb.UpgradeType_NEON:             false,
		//pb.UpgradeType_NEON_DISK:        false,
		pb.UpgradeType_SUSPENSION:       false,
		pb.UpgradeType_TRANSMISSION:     false,
		pb.UpgradeType_DIFFERENTIAL:     false,
		pb.UpgradeType_AIR_FILTER:       false,
		pb.UpgradeType_INTERCOOLER:      false,
		pb.UpgradeType_PIPES:            false,
		pb.UpgradeType_INTAKE_MAINFOLD:  false,
		pb.UpgradeType_EXHAUST_MAINFOLD: false,
		pb.UpgradeType_EXHAUST_OUTLET:   false,
		pb.UpgradeType_WESTGATE:         false,
		pb.UpgradeType_ECU:              false,
		pb.UpgradeType_OIL_COOLER:       false,
		pb.UpgradeType_OIL_INJECTORS:    false,
		pb.UpgradeType_RADIATOR:         false,
	}
	gradesPrice := map[pb.UpgradeGrade]int32{
		pb.UpgradeGrade_WHITE:  5,
		pb.UpgradeGrade_GREEN:  10,
		pb.UpgradeGrade_BLUE:   20,
		pb.UpgradeGrade_VIOLET: 40,
		pb.UpgradeGrade_YELLOW: 80,
		pb.UpgradeGrade_ORANGE: 160,
		pb.UpgradeGrade_RED:    500,
	}
	var baseId int32 = 100
	for upgradeType, hasRedGrade := range upgradeTypes {
		for grade, price := range gradesPrice {
			if grade == pb.UpgradeGrade_RED && !hasRedGrade {
				continue
			}
			blueprint := createBlueprint(upgradeType, grade, price, baseId)
			result.Items = append(result.Items, blueprint)
			baseId++
		}
		baseId += 100
	}

	return result
}

func createBlueprint(upgradeType pb.UpgradeType, grade pb.UpgradeGrade, price int32, baseId int32) *pb.BaseBlueprint {
	icon := "blueprint_" + upgradeType.String() + "_" + grade.String()
	return &pb.BaseBlueprint{
		Base: &pb.BaseItem{
			BaseId: baseId,
			Type:   pb.ItemType_BLUEPRINT,
			Price: &pb.Money{
				UpgradePoints: price,
			},
			Icon:      strings.ToLower(icon),
			ShopIndex: 1,
		},
		UpgradeType: upgradeType,
		Grade:       grade,
	}
}

func FindCar(baseId int32) *pb.BaseCar {
	for _, item := range Database.CarDatabase.Items {
		if item.BaseId == baseId {
			return item
		}
	}
	return nil
}

func FindColor(baseId int32) *pb.BaseColor {
	for _, item := range Database.ColorDatabase.Items {
		if item.BaseId == baseId {
			return item
		}
	}
	return nil
}

func FindDecal(baseId int32) *pb.BaseDecal {
	for _, item := range Database.DecalDatabase.Items {
		if item.BaseId == baseId {
			return item
		}
	}
	return nil
}

func FindTrack(baseId int32) *pb.BaseTrack {
	for _, item := range Database.TrackDatabase.Items {
		if item.BaseId == baseId {
			return item
		}
	}
	return nil
}

func FindBlueprint(baseId int32) *pb.BaseBlueprint {
	for _, item := range Database.BlueprintDatabase.Items {
		if item.Base.BaseId == baseId {
			return item
		}
	}
	return nil
}

func FindTools(baseId int32) *pb.BaseTools {
	for _, item := range Database.ToolsDatabase.Items {
		if item.Base.BaseId == baseId {
			return item
		}
	}
	return nil
}
