package protohandler

import (
	"context"
	"log"
	pb "order-service/gen/proto"
	"time"
)

func HitServer(client pb.ProductServiceClient, productId string) (*pb.ProductOrderResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	req := &pb.ProductOrderRequest{ProductId: productId}

	resp, err := client.GetProductOrderDetail(ctx, req)

	if err != nil {
		log.Println(err)
		return &pb.ProductOrderResponse{}, err
	}
	log.Println(resp)

	return resp, nil

}
