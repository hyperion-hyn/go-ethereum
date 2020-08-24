package eth

import "github.com/ethereum/go-ethereum/core"

func (b *EthAPIBackend) ChainContext() core.ChainContext {
	return b.eth.blockchain
}
