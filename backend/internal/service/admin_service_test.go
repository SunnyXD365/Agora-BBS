package service

import (
	"context"
	"errors"
	"testing"

	"agora-backend/internal/model"
)

func validAdminUserUpdate() *model.UpdateAdminUserReq {
	return &model.UpdateAdminUserReq{
		Username: "admin_user", Email: "admin@example.test", Role: "admin", Status: "active",
		UnlockLevel: 3, TrustScore: 20, Reason: "测试管理员调整保护规则",
	}
}

func TestUpdateUserProtectsCurrentAdministrator(t *testing.T) {
	adminService := &AdminService{}
	req := validAdminUserUpdate()
	req.Role = "user"
	if _, err := adminService.UpdateUser(context.Background(), 7, 7, req); !errors.Is(err, ErrSelfDemote) {
		t.Fatalf("expected ErrSelfDemote, got %v", err)
	}

	req = validAdminUserUpdate()
	req.Status = "suspended"
	if _, err := adminService.UpdateUser(context.Background(), 7, 7, req); !errors.Is(err, ErrSelfSuspend) {
		t.Fatalf("expected ErrSelfSuspend, got %v", err)
	}
}

func TestUpdateUserRequiresEmailForAdministrator(t *testing.T) {
	adminService := &AdminService{}
	req := validAdminUserUpdate()
	req.Email = ""
	if _, err := adminService.UpdateUser(context.Background(), 1, 2, req); err == nil {
		t.Fatal("expected administrator email validation error")
	}
}
