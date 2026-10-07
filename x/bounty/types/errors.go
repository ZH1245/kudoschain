package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/bounty module sentinel errors
var (
	ErrInvalidSigner = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrBountyClaimed = errors.Register(ModuleName, 1101, "bounty already claimed")
	ErrBountyExpired = errors.Register(ModuleName, 1102, "Bounty expired")
)
