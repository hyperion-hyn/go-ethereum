// Code generated by github.com/fjl/gencodec. DO NOT EDIT.

package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var _ = (*delegationMarshaling)(nil)

// MarshalJSON marshals as JSON.
func (d Delegation) MarshalJSON() ([]byte, error) {
	type Delegation struct {
		DelegatorAddress common.UnprefixedAddress `json:"delegator_address" gencodec:"required"`
		Amount           *math.HexOrDecimal256    `json:"amount" gencodec:"required"`
		Reward           *big.Int                 `json:"reward"`
	}
	var enc Delegation
	enc.DelegatorAddress = common.UnprefixedAddress(d.DelegatorAddress)
	enc.Amount = (*math.HexOrDecimal256)(d.Amount)
	enc.Reward = d.Reward
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (d *Delegation) UnmarshalJSON(input []byte) error {
	type Delegation struct {
		DelegatorAddress *common.UnprefixedAddress `json:"delegator_address" gencodec:"required"`
		Amount           *math.HexOrDecimal256     `json:"amount" gencodec:"required"`
		Reward           *big.Int                  `json:"reward"`
	}
	var dec Delegation
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.DelegatorAddress == nil {
		return errors.New("missing required field 'delegator_address' for Delegation")
	}
	d.DelegatorAddress = common.Address(*dec.DelegatorAddress)
	if dec.Amount == nil {
		return errors.New("missing required field 'amount' for Delegation")
	}
	d.Amount = (*big.Int)(dec.Amount)
	if dec.Reward != nil {
		d.Reward = dec.Reward
	}
	return nil
}