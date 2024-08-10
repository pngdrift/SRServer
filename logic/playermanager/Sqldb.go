package playermanager

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"srserver/conf"
	"srserver/logic/chat"
	"srserver/logic/top"
	"srserver/logic/tournaments"
	"srserver/logic/utils"
	pb "srserver/proto/out"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/proto"
)

var sqlDb *sql.DB

func init() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		conf.DATABASE_USERNAME,
		conf.DATABASE_PASSWORD,
		conf.DATABASE_ADDRESS,
		conf.DATABASE_NAME)

	sqlDb, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = sqlDb.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS players (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		uid VARCHAR(255) UNIQUE,
		userProto BLOB
	)`
	_, err = sqlDb.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create players table: %v", err)
	}
	log.Println("SQLDB Connected")

	//sqlDb.Exec("TRUNCATE TABLE players")
}

func GetDb() *sql.DB {
	return sqlDb
}

func Init(id string, conn net.Conn) *Player {
	player, err := LoadPlayer(id)
	if err != nil {
		log.Fatalf("Failed to load player: %v", err)
	}
	if player == nil {
		player = &Player{
			User: utils.GetDefaultUser(),
			Uid:  id,
		}
		err = player.Save()
		if err != nil {
			log.Fatalf("Failed to save new player: %v", err)
		}
		return Init(id, conn) // Kostil
	}
	player.Conn = conn
	player.User.Chat = chat.GetInitData()
	player.User.Tournaments = tournaments.GetInitData()
	player.User.Top = top.GetInitData(false)
	player.User.PointsEnemies.IsNeedUpdate = true
	player.User.Enemies.CarId = 0
	player.CheckResetTime()
	return player
}

func LoadPlayer(uid string) (*Player, error) {
	query := `SELECT id, userProto FROM players WHERE uid = ?`
	row := GetDb().QueryRow(query, uid)

	var player Player
	var userData []byte
	err := row.Scan(&player.Id, &userData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Player not found
		}
		return nil, err
	}

	player.Uid = uid
	player.User = &pb.User{}
	err = proto.Unmarshal(userData, player.User)
	if err != nil {
		return nil, err
	}
	player.User.Id = player.Id
	player.User.Info.Id = player.Id
	return &player, nil
}

func FindPlayerById(id int64) *Player {
	query := `SELECT id, uid, userProto FROM players WHERE id = ?`
	row := GetDb().QueryRow(query, id)
	var player Player
	var userData []byte
	err := row.Scan(&player.Id, &player.Uid, &userData)
	if err != nil {
		//log.Println("FindPlayerById err", err)
		return nil
	}
	player.User = &pb.User{}
	proto.Unmarshal(userData, player.User)
	//player.Conn = GetConnByPlayer(player)

	return &player
}

func RandPlayer() *Player {
	query := `
	SELECT id, uid, userProto FROM players
	ORDER BY RAND()
	LIMIT 1
	`
	row := GetDb().QueryRow(query)
	var player Player
	var userData []byte
	err := row.Scan(&player.Id, &player.Uid, &userData)
	if err != nil {
		log.Println("RandPlayer err", err)
		return nil
	}
	player.User = &pb.User{}
	proto.Unmarshal(userData, player.User)
	return &player
}

var BannedPlayers = make(map[string]bool)
