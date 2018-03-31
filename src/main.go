package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/whaangbuu/go-grpc-server/pb"
	emp "github.com/whaangbuu/go-grpc-server/src/server"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		fmt.Printf("Metadata Received: %v\n", md)
	}

	for _, e := range emp.Employees {
		if req.BadgeNumber == e.BadgeNumber {
			return &pb.EmployeeResponse{Employee: &e}, nil
		}
	}

	return nil, errors.New("Employee not found")
}

func (s *employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {
	for _, e := range emp.Employees {
		stream.Send(&pb.EmployeeResponse{Employee: &e})
	}
	return nil
}

func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {
	md, ok := metadata.FromOutgoingContext(stream.Context())

	if ok {
		fmt.Printf("Receiving photo for badge number: %v\n", md["badgenumber"][0])
	}
	imgData := []byte{}

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("File received with length: %v\n", len(imgData))
			return stream.SendAndClose(&pb.AddPhotoResponse{IsOk: true})
		}

		if err != nil {
			return err
		}

		fmt.Printf("Received %v bytes\n", len(data.Data))
		imgData = append(imgData, data.Data...)
	}
}

func (s *employeeService) Save(ctx context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	return nil, nil
}

func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	for {
		empl, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		emp.Employees = append(emp.Employees, *empl.Employee)
		stream.Send(&pb.EmployeeResponse{Employee: empl.Employee})
	}

	for _, employee := range emp.Employees {
		fmt.Println(employee)
	}
	return nil
}
