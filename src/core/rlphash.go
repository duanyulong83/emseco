package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

func RlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
/*func RlpHashOfString(str string) (h common.Hash) {
	strinterfaces := make([]interface{}, len(strs))
	for _,str := range testhash {
       strinterfaces = append(strinterfaces, str)
    }
	return rlpHash(strinterfaces)
}
*/
