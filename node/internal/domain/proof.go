package domain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 18

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	pow := &ProofOfWork{
		Block:  b,
		Target: target,
	}
	return pow
}

func (pow *ProofOfWork) InitData(nounce int64) (data []byte) {
	data = bytes.Join(
		[][]byte{
			[]byte(pow.Block.PrevHash),
			[]byte(pow.Block.Data),
			pow.ToHex(int64(nounce)),
			pow.ToHex(int64(Difficulty)),
		},
		[]byte{})
	return
}

func (pow *ProofOfWork) Run() (int64, string) {
	var intHash big.Int
	var hash [32]byte

	nounce := int64(0)

	for nounce < math.MaxInt64 {
		data := pow.InitData(nounce)
		hash = sha256.Sum256(data)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nounce++
		}
	}

	return nounce, bytesToHex(hash[:])
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nounce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func (pow *ProofOfWork) ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	if err := binary.Write(buff, binary.BigEndian, num); err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func bytesToHex(data []byte) string {
	return fmt.Sprintf("%x", data)
}
