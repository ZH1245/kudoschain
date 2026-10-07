package keeper

import (
	"context"
	"errors"

	"kudoschain/x/kudos/types"

	"cosmossdk.io/collections"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ShowKudo(ctx context.Context, req *types.QueryShowKudoRequest) (*types.QueryShowKudoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	kudo, err := q.k.Kudos.Get(ctx, req.Id)
	if errors.Is(err, collections.ErrNotFound) {
		return nil, status.Error(codes.NotFound, "kudo not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryShowKudoResponse{Kudo: &kudo}, nil
}
