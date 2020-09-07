package microstaking

// ATLAS add for rpc

type NodeKeyRPC = [48]Uint8

type Map3NodeRPC struct {
	Map3Address     Address      `json:"Map3Address" storage:"slot=0,offset=0"`
	OperatorAddress Address      `json:"OperatorAddress" storage:"slot=1,offset=0"`
	NodeKeys        []NodeKeyRPC `json:"NodeKeys" storage:"slot=2,offset=0"`
	Commission      Commission_  `json:"Commission" storage:"slot=3,offset=0"`
	Description     Description_ `json:"Description" storage:"slot=6,offset=0"`
	CreationHeight  BigInt       `json:"CreationHeight" storage:"slot=11,offset=0"`
	Age             Decimal      `json:"Age" storage:"slot=12,offset=0"`
	Status          Uint8        `json:"Status" storage:"slot=13,offset=0"`
	PendingEpoch    BigInt       `json:"PendingEpoch" storage:"slot=14,offset=0"`
	ActivationEpoch BigInt       `json:"ActivationEpoch" storage:"slot=14,offset=0"`
	ReleaseEpoch    Decimal      `json:"ReleaseEpoch" storage:"slot=15,offset=0"`
}

type Map3NodeWrapperRPC struct {
	Map3Node               Map3NodeRPC        `json:"Map3Node" storage:"slot=0,offset=0"`
	Microdelegations       []Microdelegation_ `json:"Microdelegations" storage:"slot=16,offset=0"`
	RedelegationReference  Address            `json:"RedelegationReference" storage:"slot=18,offset=0"`
	AccumulatedReward      BigInt             `json:"AccumulatedReward" storage:"slot=19,offset=0"`
	TotalDelegation        BigInt             `json:"TotalDelegation" storage:"slot=20,offset=0"`
	TotalPendingDelegation BigInt             `json:"TotalPendingDelegation" storage:"slot=21,offset=0"`
}
