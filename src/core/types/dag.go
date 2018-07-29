package types

import (
    "container/list"
    "sync"
    "github.com/ethereum/go-ethereum/common"
  )

type Anticone struct {
  AntiValue int
  Queue *list.List
}

type BlockDag struct {
	Tips      map[common.Hash]*Header
	Graph     map[common.Hash]*Header
	BlueGraph map[common.Hash]*Header
	Antic     []*Anticone
	// bc           *BlockChain
	VisitTime    int // used for score
	Coinbase     common.Address
	GenesisBlock *Header
	Lock         *sync.RWMutex
}

type DagChain struct {
	DagGraph   map[common.Hash]*Header
	OneChain []*Header
	Time     int
}

