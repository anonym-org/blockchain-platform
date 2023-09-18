package domain

import (
	"encoding/json"
	"log"
)

type Blockchain struct {
	CurrentHash string
}

type Block struct {
	Hash     string `json:"hash"`
	Data     string `json:"data"`
	PrevHash string `json:"prev_hash"`
	Nounce   int    `json:"nounce"`
}

type BlockDTO struct {
	Data string `json:"data"`
}

func NewBlock(data string, prevHash string) *Block {
	block := &Block{
		Hash:     "",
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
	return NewBlock("Genesis", "")
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
