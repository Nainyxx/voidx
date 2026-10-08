package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"example.com/m/models"
	"example.com/m/repository"
	userservice "example.com/m/userService"
	"example.com/m/utils"
	userv1 "example.com/voidx/proto/UserService"
)

// Интерфейс описан там, где используется: хендлеру нужно только это.
type UserUsecase interface {
	GetUserProfileByID(ctx context.Context, id uuid.UUID) (*models.UserProfile, error)
	CreateUserProfile(ctx context.Context, id uuid.UUID, username, name, surname, phone string) (*models.UserProfile, error)
	UpdateUserProfile(ctx context.Context, in *userservice.UpdateUserProfileInput) (*models.UserProfile, error)
	DeleteUserProfile(ctx context.Context, id uuid.UUID) error
}

type GrpcHandler struct {
	userv1.UnimplementedUserServiceServer
	svc UserUsecase
}

func NewGrpcHandler(svc UserUsecase) *GrpcHandler {
	return &GrpcHandler{svc: svc}
}

func (h *GrpcHandler) GetUserProfileByID(ctx context.Context, req *userv1.GetUserProfileByIdRequest) (*userv1.GetUserProfileByIdResponse, error) {
	id, err := parseID(req.GetUserId())
	if err != nil {
		return nil, err
	}

	p, err := h.svc.GetUserProfileByID(ctx, id)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &userv1.GetUserProfileByIdResponse{User: toProtoUser(p)}, nil
}

func (h *GrpcHandler) CreateUserProfile(ctx context.Context, req *userv1.CreateUserProfileRequest) (*userv1.CreateUserProfileResponse, error) {
	id, err := parseID(req.GetUserId())
	if err != nil {
		return nil, err
	}

	p, err := h.svc.CreateUserProfile(ctx, id, req.GetUsername(), req.GetName(), req.GetSurname(), req.GetPhone())
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &userv1.CreateUserProfileResponse{User: toProtoUser(p)}, nil
}

func (h *GrpcHandler) UpdateUserProfile(ctx context.Context, req *userv1.UpdateUserProfileRequest) (*userv1.UpdateUserProfileResponse, error) {
	id, err := parseID(req.GetUserId())
	if err != nil {
		return nil, err
	}
	in := &userservice.UpdateUserProfileInput{
		UserID:         id,
		Username:       req.Username,
		Name:           req.Name,
		Surname:        req.Surname,
		Phone:          req.Phone,
		Description:    req.Description,
		AvatarImageURL: req.AvatarImageUrl,
	}

	p, err := h.svc.UpdateUserProfile(ctx, in)
	if err != nil {
		return nil, toGRPCError(err)
	}
	return &userv1.UpdateUserProfileResponse{User: toProtoUser(p)}, nil
}

func (h *GrpcHandler) DeleteUserProfile(ctx context.Context, req *userv1.DeleteUserProfileRequest) (*emptypb.Empty, error) {
	id, err := parseID(req.GetUserId())
	if err != nil {
		return nil, err
	}

	if err := h.svc.DeleteUserProfile(ctx, id); err != nil {
		return nil, toGRPCError(err)
	}
	return &emptypb.Empty{}, nil
}

// helpers

func parseID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	return id, nil
}

func toProtoUser(p *models.UserProfile) *userv1.User {
	return &userv1.User{
		UserId:         p.UserID.String(),
		Username:       p.Username,
		Name:           p.Name,
		Surname:        p.Surname,
		Phone:          p.Phone,
		Description:    p.Description,
		AvatarImageUrl: p.AvatarImageURL,
		CreatedAt:      timestamppb.New(p.CreatedAt),
		UpdatedAt:      timestamppb.New(p.UpdatedAt),
	}
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, userservice.ErrInvalidUserID),
		errors.Is(err, userservice.ErrNilUpdate):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, utils.ErrValidation):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, repository.ErrNotFound):
		return status.Error(codes.NotFound, "user not found")

	case errors.Is(err, repository.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, "user already exists")

	default:
		slog.Error("internal error", "err", err)
		return status.Error(codes.Internal, "internal error")
	}
}
