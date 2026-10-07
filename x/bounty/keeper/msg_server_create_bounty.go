package keeper

import (
	"context"
	"strconv"

	"kudoschain/x/bounty/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateBounty(ctx context.Context, msg *types.MsgCreateBounty) (*types.MsgCreateBountyResponse, error) {
	creator, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

	if err := msg.Reward.Validate(); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, err.Error())
	}
	if !msg.Reward.IsPositive() {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "reward must be greater than zero")
	}

	// TODO: Handle the message
	if error := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(creator), types.ModuleName, sdk.NewCoins(msg.Reward)); error != nil {
		return nil, error
	}
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	id, err := k.BountySeq.Next(ctx)
	if err != nil {
		return nil, err
	}

	height := uint64(sdk.UnwrapSDKContext(ctx).BlockHeight())

	bounty := types.Bounty{
		Title:     msg.Title,
		Reward:    msg.Reward,
		Creator:   msg.Creator,
		Claimed:   false,
		ExpiresAt: height + params.BountyLifetimeBlocks,
	}
	if err := k.Bounty.Set(ctx, id, bounty); err != nil {
		return nil, err
	}
	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(sdk.NewEvent("bounty_created",
		sdk.NewAttribute("id", strconv.FormatUint(id, 10)),
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("reward", msg.Reward.String()),
		sdk.NewAttribute("expires_at", strconv.FormatUint(bounty.ExpiresAt, 10)),
	))
	return &types.MsgCreateBountyResponse{Id: id}, nil
}
