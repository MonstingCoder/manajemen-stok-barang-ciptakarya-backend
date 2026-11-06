package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/config"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/middleware"
	"github.com/monstingcoder/manajemen-stok-barang-ciptakarya-backend/internal/model"
)

type userServiceStub struct {
	regCalls   int
	loginCalls int
	user       *model.User
}

func (s *userServiceStub) Reg(context.Context, *model.UserIn) (*model.User, error) {
	s.regCalls++
	return s.user, nil
}

func (s *userServiceStub) Login(context.Context, model.Email, model.Pw) (*model.User, error) {
	s.loginCalls++
	return s.user, nil
}

func (*userServiceStub) Finds(context.Context, uint, uint) ([]model.User, error) {
	return nil, nil
}

func (*userServiceStub) Update(context.Context, uint, *model.UserIn) error { return nil }
func (*userServiceStub) Delete(context.Context, uint) error                { return nil }

type tagServiceStub struct {
	deletedID uint
	foundName string
}

func (*tagServiceStub) Create(context.Context, *model.TagIn) (*model.Tag, error) {
	return &model.Tag{}, nil
}

func (*tagServiceStub) Finds(context.Context, uint, uint) ([]model.Tag, error) {
	return []model.Tag{}, nil
}

func (s *tagServiceStub) FindByName(_ context.Context, name string) ([]model.Tag, error) {
	s.foundName = name
	return nil, nil
}

func (s *tagServiceStub) Delete(_ context.Context, id uint) error {
	s.deletedID = id
	return nil
}

func TestTagDeleteUsesPathID(t *testing.T) {
	service := new(tagServiceStub)
	handler := NewTagHandler(service)
	req := httptest.NewRequest(http.MethodDelete, "/tags/42", nil)
	req.SetPathValue("id", "42")
	res := httptest.NewRecorder()

	handler.Delete(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNoContent)
	}
	if service.deletedID != 42 {
		t.Fatalf("deleted ID = %d, want 42", service.deletedID)
	}
}

func TestTagFindByNameUsesPathValue(t *testing.T) {
	service := new(tagServiceStub)
	handler := NewTagHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/tags/stationery", nil)
	req.SetPathValue("name", "stationery")
	res := httptest.NewRecorder()

	handler.FindByName(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
	}
	if service.foundName != "stationery" {
		t.Fatalf("searched name = %q, want %q", service.foundName, "stationery")
	}
}

func TestTagCreateAndFinds(t *testing.T) {
	handler := NewTagHandler(new(tagServiceStub))
	tests := []struct {
		name   string
		method func(http.ResponseWriter, *http.Request)
		req    *http.Request
		want   int
	}{
		{
			name:   "create",
			method: handler.Create,
			req:    httptest.NewRequest(http.MethodPost, "/tags", strings.NewReader(`{}`)),
			want:   http.StatusCreated,
		},
		{
			name:   "find",
			method: handler.Finds,
			req:    httptest.NewRequest(http.MethodGet, "/tags?limit=10&offset=0", nil),
			want:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			tt.method(res, tt.req)
			if res.Code != tt.want {
				t.Fatalf("status = %d, want %d; body: %s", res.Code, tt.want, res.Body)
			}
		})
	}
}

type roleServiceStub struct{}

func (*roleServiceStub) Delete(context.Context, uint) error { return nil }

func TestRoleDelete(t *testing.T) {
	handler := NewRoleHandler(new(roleServiceStub))
	req := httptest.NewRequest(http.MethodDelete, "/roles/2", nil)
	req.SetPathValue("id", "2")
	res := httptest.NewRecorder()

	handler.Delete(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNoContent)
	}
}

func TestUserRegistrationAndLoginRejectMissingFields(t *testing.T) {
	service := new(userServiceStub)
	handler := NewUserHandler(service, nil)
	tests := []struct {
		name    string
		method  func(http.ResponseWriter, *http.Request)
		body    string
		wantReg int
		wantLog int
	}{
		{
			name:    "registration",
			method:  handler.Reg,
			body:    `{"email":"user@example.com"}`,
			wantReg: 0,
			wantLog: 0,
		},
		{
			name:    "login",
			method:  handler.Login,
			body:    `{}`,
			wantReg: 0,
			wantLog: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			res := httptest.NewRecorder()

			tt.method(res, req)

			if res.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", res.Code, http.StatusBadRequest)
			}
			if service.regCalls != tt.wantReg || service.loginCalls != tt.wantLog {
				t.Fatalf("service calls = reg:%d login:%d, want reg:%d login:%d",
					service.regCalls, service.loginCalls, tt.wantReg, tt.wantLog)
			}
		})
	}
}

func TestUserHandlers(t *testing.T) {
	roleName := "member"
	service := &userServiceStub{user: &model.User{Role: &model.Role{Name: &roleName}}}
	secret := "test-secret"
	handler := NewUserHandler(service, &config.Config{
		JWTSecret:     &secret,
		JWTExpireTime: durationPtr(time.Minute),
	})
	userBody := `{"name":"Test","email":"test@example.com","password":"Secure1!","role_id":2}`

	tests := []struct {
		name   string
		method func(http.ResponseWriter, *http.Request)
		req    *http.Request
		want   int
	}{
		{
			name:   "register",
			method: handler.Reg,
			req:    httptest.NewRequest(http.MethodPost, "/reg", strings.NewReader(userBody)),
			want:   http.StatusCreated,
		},
		{
			name:   "login",
			method: handler.Login,
			req:    httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"email":"test@example.com","password":"Secure1!"}`)),
			want:   http.StatusOK,
		},
		{
			name:   "find users",
			method: handler.Finds,
			req:    httptest.NewRequest(http.MethodGet, "/users?limit=10&offset=0", nil),
			want:   http.StatusOK,
		},
		{
			name:   "update user",
			method: handler.Update,
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPatch, "/users/1", strings.NewReader(userBody))
				req.SetPathValue("id", "1")
				ctx := context.WithValue(req.Context(), middleware.UserIDKey{}, uint(1))
				return req.WithContext(ctx)
			}(),
			want: http.StatusOK,
		},
		{
			name:   "delete user",
			method: handler.Delete,
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
				req.SetPathValue("id", "1")
				return req
			}(),
			want: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			tt.method(res, tt.req)
			if res.Code != tt.want {
				t.Fatalf("status = %d, want %d; body: %s", res.Code, tt.want, res.Body)
			}
		})
	}
	if service.regCalls != 1 || service.loginCalls != 1 {
		t.Fatalf("service calls = reg:%d login:%d, want reg:1 login:1", service.regCalls, service.loginCalls)
	}
}

type itemServiceStub struct{}

func (*itemServiceStub) Create(context.Context, *model.ItemIn) (*model.Item, error) {
	return &model.Item{}, nil
}
func (*itemServiceStub) Finds(context.Context, uint, uint) ([]model.Item, error) {
	return []model.Item{}, nil
}
func (*itemServiceStub) FindByName(context.Context, string) ([]model.Item, error) {
	return []model.Item{}, nil
}
func (*itemServiceStub) Update(context.Context, uint, *model.ItemIn) error { return nil }
func (*itemServiceStub) Delete(context.Context, uint) error                { return nil }

func TestItemHandlers(t *testing.T) {
	handler := NewItemHandler(new(itemServiceStub))
	tests := []struct {
		name   string
		method func(http.ResponseWriter, *http.Request)
		req    *http.Request
		want   int
	}{
		{
			name:   "create",
			method: handler.Create,
			req:    httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(`{}`)),
			want:   http.StatusCreated,
		},
		{
			name:   "find",
			method: handler.Finds,
			req:    httptest.NewRequest(http.MethodGet, "/items?limit=10&offset=0", nil),
			want:   http.StatusOK,
		},
		{
			name:   "find by name",
			method: handler.FindByName,
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/items/pencil", nil)
				req.SetPathValue("name", "pencil")
				return req
			}(),
			want: http.StatusOK,
		},
		{
			name:   "update",
			method: handler.Update,
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPatch, "/items/1", strings.NewReader(`{}`))
				req.SetPathValue("id", "1")
				return req
			}(),
			want: http.StatusOK,
		},
		{
			name:   "delete",
			method: handler.Delete,
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodDelete, "/items/1", nil)
				req.SetPathValue("id", "1")
				return req
			}(),
			want: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := httptest.NewRecorder()
			tt.method(res, tt.req)
			if res.Code != tt.want {
				t.Fatalf("status = %d, want %d; body: %s", res.Code, tt.want, res.Body)
			}
		})
	}
}

func durationPtr(d time.Duration) *time.Duration { return &d }
