// Stocks
package loms

import (
	"context"
	"route256/loms/internal/converter/server"
	"route256/loms/internal/model"
	"route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Stocks controller
func (s *Server) Stocks(ctx context.Context, req *loms_v1.StocksRequest) (*loms_v1.StocksResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	stocks, err := s.service.Stocks(ctx, model.SKU(req.GetSku()))
	if err != nil {
		return &loms_v1.StocksResponse{}, status.Errorf(codes.Internal, err.Error())
	}

	res := server.StocksToRes(stocks)

	return &res, nil
}
