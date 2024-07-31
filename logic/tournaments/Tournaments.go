package tournaments

import (
	pb "srserver/proto/out"
)

func GetInitData() *pb.UserTournaments {
	return &pb.UserTournaments{
		ActiveTournaments:   []*pb.UserTournament{},
		FinishedTournaments: []*pb.Tournament{},
	}
}
