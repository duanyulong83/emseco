package core

import (
	"fmt"
  "testing"
  "core/types"
  "github.com/ethereum/go-ethereum/common"
  )

func TestDagChain(t *testing.T) {
	fmt.Println("TestDagChain start!")
  var b *types.Header
  var genesis types.Header
  var graph map[common.Hash] *types.Header

  graph = make(map[common.Hash] *types.Header)
  
  var nodes = []attrs{
                     attrs{"B", []string{"genesis"}},
                     attrs{"C", []string{"genesis"}},
                     attrs{"D", []string{"genesis"}},
                     attrs{"E", []string{"genesis"}},
                     attrs{"F", []string{"B", "C"}},
                     attrs{"H", []string{"C", "D", "E"}},
                     attrs{"I", []string{"E"}},
                     attrs{"J", []string{"F", "H"}},
                     attrs{"K", []string{"B", "H", "I"}},
                     attrs{"M", []string{"F", "K"}},
                     attrs{"L", []string{"D", "I"}},
  }

  /*
   * initialize the genesis
   */

  genesis.Name = "genesis"
  genesis.ParentsHash = []common.Hash{}
  genesis.Score = 0
  genesis.ThisHash=RlpHash(genesis.Name);
  graph[genesis.ThisHash] = &genesis

  for _,x := range nodes {
    b = new(types.Header)
    b.ParentsHash = []common.Hash{}
    //b.references = make([]Address, 0, 8)

    b.Name = x.name
    b.ThisHash=RlpHash(b.Name);
    b.InLevel = -1

    for _, name := range x.parents {
       b.ParentsHash = append(b.ParentsHash, RlpHash(name))
    }
    graph[b.ThisHash] = b
  }

  vc := NewDagChain(graph)

  vc.DFS(vc.DagGraph[RlpHash("V")])
  vc.DFS_output()
  fmt.Println("TestDagChain finished!")
}

