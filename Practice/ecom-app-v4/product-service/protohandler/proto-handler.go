package protohandler

import (
	pb "product-service/gen/proto"
	"product-service/internal/products"
)

type ProtoHandler struct {
	prodConf *products.Conf
	pb.UnimplementedProductServiceServer
}

func NewProtoHandler(p *products.Conf) *ProtoHandler {
	return &ProtoHandler{
		prodConf: p,
	}
}
