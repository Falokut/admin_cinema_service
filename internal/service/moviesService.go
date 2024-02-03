package service

import (
	"context"
	"fmt"
	"strings"

	movies_service "github.com/Falokut/admin_movies_service/pkg/admin_movies_service/v1/protos"
	"google.golang.org/grpc"
)

type moviesService struct {
	client movies_service.MoviesServiceV1Client
}

func NewMoviesService(conn *grpc.ClientConn) *moviesService {
	return &moviesService{
		client: movies_service.NewMoviesServiceV1Client(conn),
	}
}

// returns movie duration in minutes if movie exist, othervise error
func (s *moviesService) GetMoviesDuration(ctx context.Context, ids []int32) (map[int32]uint32, error) {
	res, err := s.client.GetMoviesDuration(ctx,
		&movies_service.GetMoviesDurationRequest{MoviesIDs: convertIntSliceTOString(ids)})
	if err != nil {
		return map[int32]uint32{}, err
	}
	return res.Durations, nil
}

func convertIntSliceTOString(nums []int32) string {
	s := make([]string, len(nums))
	for i, num := range nums {
		s[i] = fmt.Sprint(num)
	}

	return strings.Join(s, ",")
}
