package domain

import (
	"encoding/json"
	"log"
)

type Blockchain struct {
	CurrentHash []byte
}

type Block struct {
	Hash     []byte `json:"hash"`
	Data     string `json:"data"`
	PrevHash []byte `json:"prev_hash"`
	Nounce   int    `json:"nounce"`
}

func NewBlock(data string, prevHash []byte) *Block {
	block := &Block{
		Hash:     []byte{},
		Data:     data,
		PrevHash: prevHash,
		Nounce:   0,
	}
	pow := NewProof(block)
	nounce, hash := pow.Run()
	block.Nounce = nounce
	block.Hash = hash

	return block
}

func Genesis() *Block {
	return NewBlock("Genesis", []byte{})
}

func (b *Block) Serialize() []byte {
	data, err := json.Marshal(b)
	if err != nil {
		log.Panic(err)
	}
	return data
}

func Deserialize(data []byte) *Block {
	var block Block
	if err := json.Unmarshal(data, &block); err != nil {
		log.Panic(err)
	}

	return &block
}
