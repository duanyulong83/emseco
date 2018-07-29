package core

import (
  "fmt"
  "container/list"
	"math/rand"
	"sort"
	"time"
	"strconv"

	"core/types"
	"github.com/ethereum/go-ethereum/common"
  )
const kValue int = 3
type Header types.Header
type BlockDag types.BlockDag
type DagChain types.DagChain

const (
	WHITE = 0
	GREY  = 1
	BLACK = 2
)

const (
  HashLength = 32
  AddressLength = 20
  MaxVertices = 1024*1024
  MaxLength = 256     // max width or deepth
)

const k_value int = 3


func (bd *BlockDag)toString(header *types.Header) string{
	var str=header.Name+","+strconv.Itoa(header.Score)+"(";
	for _,parent := range header.ParentsHash {
		str=str+bd.Graph[parent].Name+";";
	}
	str=str+")";
	return str
}

/*
 * create genesis node in BlockDag
 */
func NewBlockDag() *BlockDag {
  g := new(BlockDag)

  if g != nil {
    g.Tips = make(map[common.Hash] *types.Header)
    g.Graph = make(map[common.Hash] *types.Header)
    g.BlueGraph = make(map[common.Hash] *types.Header)

    g.Antic = make([]*types.Anticone, 0, MaxLength)

    fmt.Println("NewBlockDag success!")
  }

  return g
}
/*
 * receive a block from other nodes, add it to BlockDag
 */
func (bd *BlockDag) Insert(h *types.Header,parentHashs []common.Hash) error {
	//var score int

	/*if _, ok := bd.graph[h.TxHash]; ok {
		log.Warn("Header already exist")
		return nil
	}*/

	//log.Info("BlockDag insert header", "hash", h.TxHash, "Number", h.Number)

	// conert the h -> hin
	hin := types.CopyHeader(h)
	hin.Name=h.Name
	hin.ThisHash=RlpHash(hin.Name);

	// set the parent pointer
	/*bd.lock.Lock()
	bd.graph[h.TxHash] = hin
	bd.lock.Unlock()*/

	//set hin parents
	for _,hash := range parentHashs {
			hin.ParentsHash = append(hin.ParentsHash, hash)

			// remove from the tips
			if _, okin := bd.Tips[hash]; okin {
				delete(bd.Tips, hash)
			}
	}

	// set the score
	hin.Score = bd.GetScore(hin)
	//fmt.Println("cal score!",hin.Name,":score:",hin.Score)
	// set the tips of BlockDag
	bd.Tips[hin.ThisHash] = hin
	bd.Graph[hin.ThisHash] = hin
	fmt.Println(bd.toString(hin))
	return nil
}

/*
 * use the DF algorithm to get score value
 */
func (bd *BlockDag) GetScore(hin *types.Header) int {
	// initialize the graph
//	bd.lock.RLock()
//	defer bd.lock.RUnlock()

	bd.VisitTime = 0
	rand.Seed(time.Now().UnixNano())
	r := rand.Int()
	bd.searchScore(hin, r)

	return bd.VisitTime
}

func (bd *BlockDag) searchScore(hin *types.Header, r int) {
	hin.Color = r

	for _,hash := range hin.ParentsHash {
		if p, ok := bd.Graph[hash]; ok {
			if p.Color != r {
				bd.searchScore(p, r)
				bd.VisitTime += 1
			}
		}
	}

}

/*
 * get blue-set of DAG
 */
func (g *BlockDag)BlueCalc() error {
  var step common.Hash
  var blueChain []*types.Header
  var isStart bool = true
  var loop []common.Hash

  blueChain = make([]*types.Header, 0, MaxVertices)

  for {
    if isStart == true {
      loop = make([]common.Hash, 0, MaxLength)
      for _,header := range g.Tips {
      	loop=append(loop, header.ThisHash)
      }
      isStart = false
    } else {
      loop = g.Graph[step].ParentsHash
      if len(loop) == 1 {
        //if genesis, ok := g.graph[g.GenesisBlock.TxHash]; ok {
          if g.Graph[loop[0]].Name == "genesis" {
            fmt.Println("Reached Genesis: ", g.Graph[loop[0]].Name)
            blueChain = append(blueChain, g.Graph[loop[0]])
            break
          }
        //}
      }
    }

    for i:=0;i<len(loop);i++ {
      fmt.Println("loop[" , i , "]=" , g.Graph[loop[i]].Name)
    }

    step = loop[0]
    for _, x := range loop {
      if g.Graph[step].Score < g.Graph[x].Score {
        step = x
      }
    }

    fmt.Println(g.Graph[step].Name)
    blueChain = append(blueChain, g.Graph[step])
  }

  for i:=len(blueChain)-1;i>=0;i-- {
    x := blueChain[i]
    level := len(blueChain)- 1 -i

    g.BlueGraph[x.ThisHash] = x
    x.InLevel = level

    if level == 0 {
      g.anticone_process(0, x, nil)
    }

    fmt.Println("x=", x.Name, "len=", len(g.BlueGraph))

    if len(x.ParentsHash) > 0 {
      for _, y := range x.ParentsHash {
        fmt.Println("  y =", g.Graph[y].Name)
        if g.anticone_process(level, x, g.Graph[y]) == true {
          fmt.Println("     anticone is true: ", g.Graph[y].Name)
          g.BlueGraph[g.Graph[y].ThisHash] = g.Graph[y]
          g.Graph[y].InLevel = level
        }
      }
    }
  }

  return nil
}

/*
 * process anticone of a block
 */
func NewAnticone() *types.Anticone {
  a := new(types.Anticone)
  if a != nil {
    a.AntiValue = 0
    a.Queue = list.New()
  }

  return a
}

func (g *BlockDag)anticone_process(level int, s *types.Header, c *types.Header) bool {
  var minimal int = MaxVertices
  var ok bool = false

  fmt.Println("anticone_process: level =", level)

  if len(g.Antic) == level {
    anti := NewAnticone()
    anti.Queue.PushBack(s.Name)
    g.Antic = append(g.Antic, anti)
  }

  // 1. find all the parents of c, get the one whose level is minimal
  if c != nil {
    for _, parentHash := range c.ParentsHash {
      if g.BlueGraph[parentHash]!=nil && g.BlueGraph[parentHash].InLevel >= 0 && g.BlueGraph[parentHash].InLevel < minimal {
        minimal = g.BlueGraph[parentHash].InLevel
      }
    }
  }

  // 2. judge the antiValue between the minimal and the current
  if minimal != MaxVertices {
    for i:= minimal+1;i<= level;i++ {
      anti := g.Antic[i]
      if anti.AntiValue < k_value {
        anti.AntiValue++
        ok = true
      } else {
        break
      }
    }
  }

  // 3. add current block to this step
  if ok {
    g.Antic[level].Queue.PushBack(c.Name)
    g.Antic[level].AntiValue++
  }

  return ok
}

/*
 * Deep-First-Search in Graph
 */
func NewDagChain(blueGraph map[common.Hash]*types.Header) *DagChain {
  v := new(DagChain)
  v.DagGraph = make(map[common.Hash] *types.Header)
  v.Time = 0
  v.OneChain = make([]*types.Header, 0, MaxLength)

  for _,header  := range blueGraph {
    vb := new(types.Header)

    vb.Name = header.Name
    vb.Color = WHITE
    vb.Pi = nil
    v.DagGraph[header.ThisHash] = vb
  }

  for hash, header := range v.DagGraph {
    if value, ok := blueGraph[hash]; ok {
      for _, n := range value.ParentsHash {
        header.ParentsHash = append(header.ParentsHash, n)
      }
    }
  }

  return v
}


func (vc *DagChain)DFS(u *types.Header) {
  if u == nil {
    for _,x := range vc.DagGraph {
      if x.Color == WHITE {
        vc.DFS_Visit(x)
      }
    }
  } else {
    vc.DFS_Visit(u)
  }

  return
}

func (vc *DagChain)DFS_Visit(u *types.Header) {
  vc.Time += 1
  u.DiscoveryTime = vc.Time
  u.Color = GREY

  for _,x := range(u.ParentsHash) {
    if vc.DagGraph[x].Color == WHITE {
      vc.DagGraph[x].Pi = u
      vc.DFS_Visit(vc.DagGraph[x])
    }
  }
  u.Color = BLACK
  vc.Time += 1
  u.FinishTime = vc.Time
}

func (vc *DagChain) Len() int {
  return len(vc.OneChain)
}

func (vc *DagChain) Swap(i, j int) {
  vc.OneChain[i], vc.OneChain[j] = vc.OneChain[j], vc.OneChain[i]
}

func (vc *DagChain) Less(i, j int) bool {
  return vc.OneChain[i].FinishTime < vc.OneChain[j].FinishTime
}

func (vc *DagChain)DFS_output() {
  for _, v := range vc.DagGraph {
    vc.OneChain = append(vc.OneChain, v)
  }

  fmt.Println("after sort:")
  sort.Sort(vc)
  for _, v := range vc.OneChain {
    fmt.Println("\tnode =", v.Name, "seq =", v.FinishTime)
  }
}

