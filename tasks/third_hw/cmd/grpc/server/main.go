package main

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/l-brawler-l/go_test/tasks/third_hw/accounts/models"
	"github.com/l-brawler-l/go_test/tasks/third_hw/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New() *server {
	return &server{
		accounts: make(map[string]*models.Account),
		guard:    &sync.RWMutex{},
	}
}

type server struct {
	proto.UnimplementedBankAccountsServer
	accounts map[string]*models.Account
	guard    *sync.RWMutex
}

func (s *server) Create (ctx context.Context, req *proto.CreateAccountRequest) (*proto.Empty, error) {
	if len(req.Name) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}
	s.guard.Lock()

	if _, ok := s.accounts[req.Name]; ok {
		s.guard.Unlock()

		return nil, status.Errorf(codes.AlreadyExists, "account already exists")
	}

	s.accounts[req.Name] = &models.Account{
		Name:   req.Name,
		Amount: int(req.Amount),
	}

	s.guard.Unlock()

	return &proto.Empty{}, nil
}

func (s *server) Get (ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountReply, error) {
	if len(req.Name) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}
	s.guard.RLock()

	account, ok := s.accounts[req.Name]

	s.guard.RUnlock()

	if !ok {
		return nil, status.Errorf(codes.NotFound, "account not found")
	}

	response := proto.GetAccountReply{
		Name:   account.Name,
		Amount: int32(account.Amount),
	}

	return &response, nil
}

func (s *server) Patch (ctx context.Context, req *proto.PatchAccountRequest) (*proto.Empty, error) {
	if len(req.Name) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}
	s.guard.Lock()

	if _, ok := s.accounts[req.Name]; !ok {
		s.guard.Unlock()

		return nil, status.Errorf(codes.NotFound, "account not found")
	}

	s.accounts[req.Name].Amount = int(req.Amount)

	s.guard.Unlock()

	return &proto.Empty{}, nil
}

func (s *server) Delete (ctx context.Context, req *proto.DeleteAccountRequest) (*proto.Empty, error) {
	if len(req.Name) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}
	s.guard.Lock()

	if _, ok := s.accounts[req.Name]; !ok {
		s.guard.Unlock()

		return nil, status.Errorf(codes.NotFound, "account not found")
	}

	delete(s.accounts, req.Name)

	s.guard.Unlock()

	return &proto.Empty{}, nil
}

func (s *server) Change (ctx context.Context, req *proto.ChangeAccountRequest) (*proto.Empty, error) {
	if len(req.Name) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty name")
	}
	if len(req.NewName) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "empty new name")
	}
	s.guard.Lock()

	if _, ok := s.accounts[req.Name]; !ok {
		s.guard.Unlock()
		return nil, status.Errorf(codes.NotFound, "account with that name is not exist")
	}

	if _, ok := s.accounts[req.NewName]; ok {
		s.guard.Unlock()
		return nil, status.Errorf(codes.NotFound, "account with new name already exists")
	}
	s.accounts[req.NewName] = s.accounts[req.Name]
	s.accounts[req.NewName].Name = req.NewName
	delete(s.accounts, req.Name)

	s.guard.Unlock()

	return &proto.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 4567))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	h := New()
	proto.RegisterBankAccountsServer(s, h)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}

}