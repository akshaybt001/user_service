package service

import (
	"context"
	"errors"
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

func (u *UserService) AdminLogin(ctx context.Context,req *pb.LoginRequest)(*pb.UserResponse,error){
	res,err:=u.Adapter.AdminLogin(req.Email)
	if err!=nil{
		return nil,err
	}

	if err:=helper.VerifyPassword(res.Password,req.Password);err!=nil{
		return nil,err
	}

	return &pb.UserResponse{
		Id: uint32(res.Id),
		Name: res.Name,
		Email: res.Email,
	},nil
}

func (u *UserService) SupAdminLogin(ctx context.Context,req *pb.LoginRequest)(*pb.UserResponse,error){
	res,err:=u.Adapter.SupAdminLogin(req.Email)
	if err!=nil{
		return nil,err
	}

	if err:=helper.VerifyPassword(res.Password,req.Password);err!=nil{
		return nil,err
	}
	return &pb.UserResponse{
		Id: uint32(res.Id),
		Name: res.Name,
		Email: res.Email,
	},nil
}

func (u *UserService) AddAdmin(ctx context.Context,req *pb.UserSignUpRequest)(*pb.UserResponse,error){

	if req.Name == ""{
		return nil,errors.New("the name cannot be empty")
	}
	if req.Email == ""{
		return nil,errors.New("the email cannot be empty")
	}
	if req.Password == ""{
		return nil,errors.New("the password cannot be empty")
	}

	password,err:=helper.HashPassword(req.Password)
	if err!=nil{
		fmt.Println(err.Error())
		return nil,err
	}

	reqq:=entities.Admin{
		Name: req.Name,
		Email: req.Email,
		Password: password,
	}

	userRes,err:=u.Adapter.AddAdmin(reqq)
	if err!=nil{
		return nil,err
	}
	return &pb.UserResponse{
		Id: uint32(userRes.Id),
		Name: userRes.Name,
		Email: userRes.Email,
	},nil
}
func (admin *UserService) GetAllUsers(em *pb.NoPara,srv pb.UserService_GetAllUsersServer)error{
	span:=Tracer.StartSpan("get all users")
	defer span.Finish()
	users,err:=admin.Adapter.GetAllUsers()
	if err!=nil{
		return err
	}
	for _,user:=range users{
		if err= srv.Send(&pb.UserResponse{
			Id: uint32(user.Id),
			Name: user.Name,
			Email: user.Email,
		});err!=nil{
			return err
		}
	}
	return nil
}

func (sup *UserService) GetAllAdmins(em *pb.NoPara,srv pb.UserService_GetAllAdminsServer)error{
	span:=Tracer.StartSpan("get all admins grpc")
	defer span.Finish()
	admins,err:=sup.Adapter.GetAllAdmins()
	if err!=nil{
		return err
	}
	for _,admin:=range admins{
		if err=srv.Send(&pb.UserResponse{
			Id: uint32(admin.Id),
			Name: admin.Name,
			Email: admin.Email,
		});err!=nil{
			return err
		}
	}
	return nil
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
