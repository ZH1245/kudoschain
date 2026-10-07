package keeper

import (
	"context"
	"errors"

	"kudoschain/x/bounty/types"

	"cosmossdk.io/collections"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (q queryServer) ShowBounty(ctx context.Context, req *types.QueryShowBountyRequest) (*types.QueryShowBountyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	existingBounty, err := q.k.Bounty.Get(ctx, req.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "Bounty not found")
		}
		return nil, status.Error(codes.Internal, "Internal Error")
	}
	// TODO: Process the query

	return &types.QueryShowBountyResponse{Bounty: &existingBounty}, nil
}
