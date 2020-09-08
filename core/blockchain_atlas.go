package core

import (
	"github.com/ethereum/go-ethereum/ethdb"
)

// ChainDb returns the database
func (bc *BlockChain) ChainDb() ethdb.Database { return bc.db }
