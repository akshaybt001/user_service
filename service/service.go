package service

import (
	"context"
	"fmt"

	"github.com/akshaybt001/proto_files/pb"
	"github.com/akshaybt001/user_service/adapter"
	"github.com/akshaybt001/user_service/entities"
	"github.com/akshaybt001/user_service/helper"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var (
	Tracer opentracing.Tracer
)

func RetrieveTracer(tr opentracing.Tracer) {
	Tracer = tr
}

type UserService struct {
	Adapter adapter.UserInterface
	pb.UnimplementedUserServiceServer
}

func NewUserService(adapter adapter.UserInterface) *UserService {
	return &UserService{
		Adapter: adapter,
	}
}

func (u *UserService) UserSignUp(ctx context.Context, req *pb.UserSignUpRequest) (*pb.UserResponse, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name cant be empty")
	} else if req.Email == "" {
		return nil, fmt.Errorf("email cant be empty")
	} else if req.Password == "" {
		return nil, fmt.Errorf("password cant be empty")
	}

	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	req.Password = hashed

	ogreq := entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := u.Adapter.UserSignUp(ogreq)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &pb.UserResponse{
		Id:    uint32(res.Id),
		Name:  res.Name,
		Email: res.Email,
	}, nil
}

func (u *UserService) UserLogin(ctc context.Context, req *pb.LoginRequest) (*pb.UserResponse, error) {
	res, err := u.Adapter.UserLogin(req.Email)
	if err != nil {
		return nil, fmt.Errorf("there is no such user")
	}

	if err := helper.VerifyPassword(res.Password, req.Password); err != nil {
		return nil, fmt.Errorf("wrong password")
	}

	return &pb.UserResponse{
		Id:    uint32(res.Id),
		Name:  res.Name,
		Email: res.Email,
	}, nil
}

type HealthChecker struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *HealthChecker) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *HealthChecker) Watch(in *grpc_health_v1.HealthCheckRequest, srv grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watching is not supported")
}
