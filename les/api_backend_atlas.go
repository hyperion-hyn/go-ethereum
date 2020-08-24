package les

import (
	"github.com/ethereum/go-ethereum/core"
)

func (b *LesApiBackend) ChainContext() core.ChainContext {
	return b.eth.blockchain
}