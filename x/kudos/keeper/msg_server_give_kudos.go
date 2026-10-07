package keeper

import (
	"context"

	"kudoschain/x/kudos/types"

	errorsmod "cosmossdk.io/errors"
)

func (k msgServer) GiveKudos(ctx context.Context, msg *types.MsgGiveKudos) (*types.MsgGiveKudosResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid giver address")
	}
	if _, err := k.addressCodec.StringToBytes(msg.Receiver); err != nil {
		return nil, errorsmod.Wrap(err, "invalid receiver address")
	}

	id, err := k.KudosSeq.Next(ctx)
	if err != nil {
		return nil, err
	}

	kudo := types.Kudo{Giver: msg.Creator, Receiver: msg.Receiver, Message: msg.Message}
	if err := k.Kudos.Set(ctx, id, kudo); err != nil {
		return nil, err
	}

	return &types.MsgGiveKudosResponse{Id: id}, nil
}
