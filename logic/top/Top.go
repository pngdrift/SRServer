package top

import (
	pb "srserver/proto/out"
)

func GetInitData(loadTimeTop bool) *pb.Top {
	return &pb.Top{
		PointsTop: []*pb.TopItem{},
		RatingTop: []*pb.TopItem{},
		Place:     0,
		Rating:    0,
		MinRating: 0,
		MaxRating: 0,
		TimeTop:   []*pb.TopItem{},
		TimeTopA:  []*pb.TopItem{},
		TimeTopB:  []*pb.TopItem{},
		TimeTopC:  []*pb.TopItem{},
		TimeTopD:  []*pb.TopItem{},
		TimeTopE:  []*pb.TopItem{},
		TimeTopF:  []*pb.TopItem{},
		TimeTopG:  []*pb.TopItem{},
		TimeTopH:  []*pb.TopItem{},
		TimeTopI:  []*pb.TopItem{},
		TimeTopJ:  []*pb.TopItem{},
		TimeTopK:  []*pb.TopItem{},
		TimeTopL:  []*pb.TopItem{},
	}
}
