package service

import (
	"auth-haven/internal/domain/role"
	"auth-haven/internal/domain/tenant"
	"auth-haven/internal/domain/user"
	"auth-haven/internal/utils"
	proto "auth-haven/pkg/proto"
	common "auth-haven/pkg/proto/common"
	"context"
	"errors"
	"time"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
	UserRepo   user.UserRepository
	TenantRepo tenant.TenantRepository
	RoleRepo   role.RoleRepository
}

// CreatePersonalUser creates a simple user without a company
func (s *UserService) CreatePersonalUser(ctx context.Context, req *proto.CreatePersonalUserRequest) (*common.Tokens, error) {
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		return nil, errors.New("missing required fields")
	}
	hashed, err := utils.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	newUser := &user.User{
		Email:        req.Email,
		PasswordHash: hashed, // Ideally, hash it
		FullName:     req.FullName,
	}

	createdUser, err := s.UserRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// Here you would generate JWT tokens
	tokens := &common.Tokens{
		AccessToken:  "mock-access-token-" + createdUser.ID,
		RefreshToken: "mock-refresh-token-" + createdUser.ID,
	}

	return tokens, nil
}

// CreateCompanyAndOwner creates a tenant, owner role, and owner user
func (s *UserService) CreateCompanyAndOwner(ctx context.Context, req *proto.CreateCompanyAndOwnerRequest) (*common.Tokens, error) {
	if req.CompanyName == "" || req.CompanyDomain == "" || req.OwnerEmail == "" || req.OwnerPassword == "" || req.OwnerFullName == "" {
		return nil, errors.New("missing required fields")
	}

	// 1. Create Tenant
	newTenant := &tenant.Tenant{
		Name:   req.CompanyName,
		Domain: req.CompanyDomain,
		Status: "active",
	}
	createdTenant, err := s.TenantRepo.Create(ctx, newTenant)
	if err != nil {
		return nil, err
	}

	// 2. Create Owner Role
	ownerRole := &role.Role{
		TenantID:    createdTenant.ID,
		Name:        "owner",
		Permissions: []byte(`{"admin":true}`),
		CreatedAt:   time.Now(),
	}
	createdRole, err := s.RoleRepo.Create(ctx, ownerRole)
	if err != nil {
		return nil, err
	}
	hashed, err := utils.Hash(req.OwnerPassword)
	if err != nil {
		return nil, err
	}

	// 3. Create Owner User
	ownerUser := &user.User{
		TenantID:     createdTenant.ID,
		Email:        req.OwnerEmail,
		PasswordHash: hashed, // Ideally hashed
		FullName:     req.OwnerFullName,
		RoleId:       createdRole.ID,
	}
	createdOwner, err := s.UserRepo.Create(ctx, ownerUser)
	if err != nil {
		return nil, err
	}

	// 4. Generate Tokens
	tokens := &common.Tokens{
		AccessToken:  "mock-access-token-" + createdOwner.ID,
		RefreshToken: "mock-refresh-token-" + createdOwner.ID,
	}

	return tokens, nil
}
