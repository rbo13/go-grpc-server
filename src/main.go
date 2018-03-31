package main

import (
	"log"
	"net"

	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/whaangbuu/go-grpc-server/pb"
)

const port = ":9000"

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}
	server := grpc.NewServer(opts...)
	pb.RegisterEmployeeServiceServer(server, new(employeeService))
	log.Println("Starting on port " + port)
	server.Serve(listen)
}

type employeeService struct{}

func (s *employeeService) GetByBadgeNumber(ctx context.Context, req *pb.GetByBadgeNumberRequest) (*pb.EmployeeResponse, error) {
	return nil, nil
}

func (s *employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {
	return nil
}

func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {
	return nil
}

func (s *employeeService) Save(ctx context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	return nil, nil
}

func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	return nil
}
