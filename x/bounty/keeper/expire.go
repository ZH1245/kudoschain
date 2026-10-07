package keeper

import (
	"context"
	"kudoschain/x/bounty/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RefundExpiredBounties(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	height := uint64(sdkCtx.BlockHeight())
	var expired []uint64
	error := k.Bounty.Walk(ctx, nil, func(id uint64, bounty types.Bounty) (bool, error) {
		if !bounty.Claimed && bounty.ExpiresAt <= height {
			expired = append(expired, id)
		}
		return false, nil
	})
	if error != nil {
		return error
	}
	for _, id := range expired {
		bounty, err := k.Bounty.Get(ctx, id)
		if err != nil {
			return err
		}
		creator, err := k.addressCodec.StringToBytes(bounty.Creator)
		if err != nil {
			return err
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(creator), sdk.NewCoins(bounty.Reward))
		if err != nil {
			return err
		}
		if err := k.Bounty.Remove(ctx, id); err != nil {
			return err
		}
		sdkCtx.EventManager().EmitEvent(sdk.NewEvent("bounty_refunded",
			sdk.NewAttribute("id", strconv.FormatUint(id, 10)),
			sdk.NewAttribute("creator", bounty.Creator),
			sdk.NewAttribute("reward", bounty.Reward.String()),
		))
	}
	return nil
}
