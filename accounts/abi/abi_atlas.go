package abi

// PackRevert encodes revert reason as the solidity
// spec https://solidity.readthedocs.io/en/latest/control-structures.html#revert.
func PackRevert(data string) ([]byte, error) {
	typ, _ := NewType("string", "", nil)
	b, err := (Arguments{{Type: typ}}).Pack(data)
	if err != nil {
		return nil, err
	}
	result := append([]byte{}, revertSelector...)
	return append(result, b...), nil
}