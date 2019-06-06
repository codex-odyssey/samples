package main

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/wrappers"
	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/grpc-up-and-running/samples/ch03/order-service/go/order_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

var orderMap = make(map[string]pb.Order)

type server struct {
	orderMap map[string]*pb.Order
}

func (s *server) AddOrder(ctx context.Context, orderReq *pb.Order) (*wrappers.StringValue, error) {
	orderMap[orderReq.Id] = *orderReq
	return &wrapper.StringValue{Value: "Order Added: " + orderReq.Id}, nil
}

func (s *server) GetOrderStatus(ctx context.Context, in *wrapper.StringValue) (*pb.Order, error) {

	myOrder := new(pb.Order)
	myOrder.Id = "100500"
	myOrder.Name = "Sample Order"
	myOrder.Description = "Order description of 100500 Sample Order"
	return myOrder, nil
}

func (s *server) SearchOrders(req *wrappers.StringValue, stream pb.OrderManagement_SearchOrdersServer) error {
	myOrder := new(pb.Order)
	myOrder.Id = "100500"
	myOrder.Name = "Sample Search Order"
	myOrder.Description = "Order description of 100500 Sample Search Order"

	if err := stream.Send(myOrder); err != nil {
		return err
	}

	order1 := new(pb.Order)
	order1.Id = "100501"
	order1.Name = "Sample Search Order 2"
	order1.Description = "Order description of 100501 Sample Search Order 2"
	stream.Send(order1)

	return nil
}

func (s *server) UpdateOrders(stream pb.OrderManagement_UpdateOrdersServer) error {

	ordersStr := ""
	for {
		order, err := stream.Recv()
		if err == io.EOF {
			// Finished reading order stream
			return stream.SendAndClose(&wrapper.StringValue{Value: "Orders processed " + ordersStr})
		}
		// Process order

		log.Printf("Order ID ", order.Id, " - Processed!")
		ordersStr += order.Id + "\n"
		// ...
	}
}

func (s *server) ProcessOrders(stream pb.OrderManagement_ProcessOrdersServer) error {
	order, _ := stream.Recv()
	orderList := []string{"100500", "100501"}
	comb := pb.CombinedShipment{Id: "123", OrderIDList: orderList, Status: "OK"}

	log.Printf("Order ID ", order.Id, " - Processed!")

	stream.Send(&comb)
	return nil
	//return status.Errorf(codes.Unimplemented, "method ProcessOrders not implemented")
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOrderManagementServer(s, &server{})
	///	pb.RegisterOrderInfoServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
