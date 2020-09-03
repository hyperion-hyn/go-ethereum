package restaking

// ATLAS add for rpc

type SlotPubKeyRPC = [48]Uint8

type ValidatorRPC struct {
	ValidatorAddress     Address         `json:"ValidatorAddress" `
	OperatorAddresses    []Address       `json:"OperatorAddresses" `
	SlotPubKeys          []SlotPubKeyRPC `json:"SlotPubKeys" `
	LastEpochInCommittee BigInt          `json:"LastEpochInCommittee" `
	MaxTotalDelegation   BigInt          `json:"MaxTotalDelegation"`
	Status               Uint8           `json:"Status" `
	Commission           Commission_     `json:"Commission" `
	Description          Description_    `json:"Description" `
	CreationHeight       BigInt          `json:"CreationHeight" `
}

type ValidatorWrapperRPC struct {
	Validator                 ValidatorRPC    `json:"Validator" storage:"slot=0,offset=0"`
	Redelegations             []Redelegation_ `json:"Redelegations" storage:"slot=17,offset=0"`
	Counters                  Counters_       `json:"Counters" storage:"slot=19,offset=0"`
	BlockReward               BigInt          `json:"BlockReward" storage:"slot=21,offset=0"`
	TotalDelegation           BigInt          `json:"TotalDelegation" storage:"slot=22,offset=0"`
	TotalDelegationByOperator BigInt          `json:"TotalDelegationByOperator" storage:"slot=23,offset=0"`
}
