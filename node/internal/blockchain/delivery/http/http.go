package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/BakuPukul/blockchain-platform/config"
	"github.com/BakuPukul/blockchain-platform/internal/blockchain"
	"github.com/BakuPukul/blockchain-platform/internal/domain"
	httpresponse "github.com/BakuPukul/blockchain-platform/pkg/http-response"
	"github.com/BakuPukul/blockchain-platform/pkg/network"
)

type controller struct {
	config     config.Config
	blockchain *domain.Blockchain
	handler    *http.ServeMux
	service    blockchain.Usecase
	network    *network.Network
}

func NewController(config config.Config, blockchain *domain.Blockchain, handler *http.ServeMux, service blockchain.Usecase, network *network.Network) {
	c := &controller{
		config:     config,
		blockchain: blockchain,
		handler:    handler,
		service:    service,
		network:    network,
	}

	handler.HandleFunc("/api/blocks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			c.getBlocks(w, r)
		case http.MethodPost:
			c.addBlock(w, r)
		}
	})
}

func (c *controller) getBlocks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	block, err := c.service.GetBlock(ctx, c.blockchain)
	if err != nil {
		httpresponse.WriteErrorResponse(w, r, http.StatusBadRequest, "fail to get block", "invalid_param")
		return
	}

	httpresponse.WriteSuccessResponse(w, r, http.StatusOK, "get latest block from blockchain", block.Data)
}

func (c *controller) addBlock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpresponse.WriteErrorResponse(w, r, http.StatusBadRequest, "invalid params", "invalid_param")
		return
	}
	defer r.Body.Close()

	var dto domain.BlockDTO
	if err = json.Unmarshal(body, &dto); err != nil || dto.Data == "" {
		httpresponse.WriteErrorResponse(w, r, http.StatusBadRequest, "invalid params", "invalid_param")
		return
	}

	block, err := c.service.AddBlock(ctx, c.blockchain, dto.Data)
	if err != nil {
		httpresponse.WriteErrorResponse(w, r, http.StatusBadRequest, err.Error(), "invalid_param")
		return
	}

	c.network.Broadcast(block)
	httpresponse.WriteSuccessResponse(w, r, http.StatusCreated, "data added to blockchain", block.Data)
}
