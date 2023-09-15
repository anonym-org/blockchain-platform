package controller

import (
	"net/http"

	"github.com/BakuPukul/blockchain-platform/internal/blockchain"
	"github.com/BakuPukul/blockchain-platform/internal/domain"
	httpresponse "github.com/BakuPukul/blockchain-platform/pkg/http-response"
)

type controller struct {
	blockchain *domain.Blockchain
	handler    *http.ServeMux
	service    blockchain.Usecase
}

func NewController(blockchain *domain.Blockchain, handler *http.ServeMux, service blockchain.Usecase) {
	c := &controller{
		blockchain: blockchain,
		handler:    handler,
		service:    service,
	}

	handler.HandleFunc("/api/blocks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			c.getBlocks(w, r)
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
