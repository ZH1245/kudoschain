package types

import fmt "fmt"

const DefaultBountyLifetimeBlocks uint64 = 20

// NewParams creates a new Params instance.
func NewParams(bountyLifetimeBlocks uint64) Params {
	return Params{BountyLifetimeBlocks: bountyLifetimeBlocks}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultBountyLifetimeBlocks)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.BountyLifetimeBlocks == 0 {
		return fmt.Errorf("bounty_lifetime_blocks must be greater than 0")
	}
	return nil
}
