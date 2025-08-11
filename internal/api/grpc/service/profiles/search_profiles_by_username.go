package profiles

import (
	"context"

	svc "github.com/chains-lab/profiles-proto/gen/go/profile"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s Service) SearchProfilesByUsername(ctx context.Context, req *svc.SearchProfilesByUsernameRequest) (*svc.ProfilesList, error) {
	return nil, status.New(codes.Internal, "not implemented").Err()
}
