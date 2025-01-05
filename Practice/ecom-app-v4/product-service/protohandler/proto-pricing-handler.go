package protohandler

import (
	"context"
	"log/slog"
	pb "product-service/gen/proto"
	"product-service/pkg/logkey"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p ProtoHandler) GetProductOrderDetail(ctx context.Context, req *pb.ProductOrderRequest) (*pb.ProductOrderResponse, error) {

	productID := req.GetProductId()

	prodOrder, err := p.prodConf.GetStripeProductDetail(ctx, productID)

	if err != nil {
		//TODO add trace id
		slog.Error(
			"failed to get stripe price id",
			slog.Any(logkey.ERROR, err.Error()))
		return nil, status.Errorf(codes.Internal, "Failed to get stripe price id")
	}
	//slog.Info("successfully got stripe customer id for", productID)

	pbProdOrder := &pb.ProductOrderDetails{
		PriceId: prodOrder.PriceId,
		Stock:   int64(prodOrder.Stock),
	}
	return &pb.ProductOrderResponse{ProdOrder: pbProdOrder}, nil
}

