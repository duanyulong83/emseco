package core

import (
	"fmt"
  "testing"
  "core/types"
  "github.com/ethereum/go-ethereum/common"
  )
type attrs struct {
  name string
  parents []string
}
func creating_test(bd *BlockDag) {
  var genesis types.Header
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
  bd.Graph[genesis.ThisHash] = &genesis
  
  /*
   * initialize the graph
   */
  for _, node := range nodes {
    var header *types.Header
    header=new(types.Header)
    header.ThisHash=RlpHash(node.name)
    header.Name=node.name;
    var parentHashs = make([]common.Hash, 0, len(node.parents))
    for _,name := range node.parents {
       parentHashs = append(parentHashs, RlpHash(name))
    }
    bd.Insert(header,parentHashs);
  } // for loop ends

  return
}


func TestBlue(t *testing.T) {
	fmt.Println("TestBlue start!")
  var g *BlockDag
  g = NewBlockDag()
  if g == nil {
    t.Error("NewDag failed")
  }

  creating_test(g)

  g.BlueCalc()

  fmt.Print("The BlueGraph : ")
  for _,x := range g.BlueGraph {
    fmt.Print(x.Name+" ;")
  }
  fmt.Println()
  fmt.Println("TestBlue finished")
}

