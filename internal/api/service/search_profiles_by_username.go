package service

import (
	"context"
	"fmt"

	"github.com/chains-lab/profiles-proto/gen/go/svc"
	"github.com/chains-lab/profiles-svc/internal/ape"
)

func (s Service) SearchProfilesByUsername(ctx context.Context, req *svc.SearchProfilesByUsernameRequest) (*svc.ProfilesList, error) {
	return nil, ape.RaiseInternal(fmt.Errorf("inpliment me"))
}
