package domain

import (
	"encoding/json"
	"log"
)

type Blockchain struct {
	GenesisHash string
	CurrentHash string
}

type Block struct {
	Hash     string `json:"hash"`
	Data     string `json:"data"`
	PrevHash string `json:"prev_hash"`
	Nounce   int64  `json:"nounce"`
}

type GenesisDTO struct {
	ID string `json:"id"`
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
	g := struct {
		Data any `json:"data"`
	}{
		Data: GenesisDTO{
			ID: "Genesis",
		},
	}
	v, err := json.Marshal(g)
	if err != nil {
		log.Fatal(err)
	}
	return NewBlock(string(v), "")
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
