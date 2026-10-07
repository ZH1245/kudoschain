package keeper

import (
	"context"
	"errors"
	"strconv"

	"kudoschain/x/bounty/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) ClaimBounty(ctx context.Context, msg *types.MsgClaimBounty) (*types.MsgClaimBountyResponse, error) {
	claimer, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}
	existingBounty, err := k.Bounty.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrNotFound, "bounty %d not found", msg.Id)
		}
		return nil, err
	}
	if existingBounty.Claimed {
		return nil, errorsmod.Wrapf(types.ErrBountyClaimed, "Bounty %d already claimed", msg.Id)
	}
	if height := uint64(sdk.UnwrapSDKContext(ctx).BlockHeight()); height >= existingBounty.ExpiresAt {
		return nil, errorsmod.Wrapf(types.ErrBountyExpired, "Bounty %d Expired", msg.Id)
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(claimer), sdk.NewCoins(existingBounty.Reward)); err != nil {
		return nil, err
	}
	existingBounty.Claimed = true
	if err := k.Bounty.Set(ctx, msg.Id, existingBounty); err != nil {
		return nil, err
	}
	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(sdk.NewEvent("bounty_claimed",
		sdk.NewAttribute("id", strconv.FormatUint(msg.Id, 10)),
		sdk.NewAttribute("claimer", msg.Creator)))
	return &types.MsgClaimBountyResponse{Id: msg.Id, Claimed: true}, nil
}
