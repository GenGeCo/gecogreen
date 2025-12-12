package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleBuyer  UserRole = "BUYER"
	RoleSeller UserRole = "SELLER"
	RoleAdmin  UserRole = "ADMIN"
)

type UserStatus string

const (
	UserStatusPending   UserStatus = "PENDING"
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
	UserStatusBanned    UserStatus = "BANNED"
)

type User struct {
	ID               uuid.UUID  `json:"id"`
	Email            string     `json:"email"`
	PasswordHash     string     `json:"-"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Phone            string     `json:"phone,omitempty"`
	City             string     `json:"city,omitempty"`
	Roles            []UserRole `json:"roles"`
	Status           UserStatus `json:"status"`
	EmailVerified    bool       `json:"email_verified"`
	AvatarURL        string     `json:"avatar_url,omitempty"`
	StripeCustomerID *string    `json:"stripe_customer_id,omitempty"`
	StripeAccountID  *string    `json:"stripe_account_id,omitempty"`
	LastLoginAt      *time.Time `json:"last_login_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (u *User) HasRole(role UserRole) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (u *User) IsSeller() bool { return u.HasRole(RoleSeller) }
func (u *User) IsAdmin() bool  { return u.HasRole(RoleAdmin) }
func (u *User) IsActive() bool { return u.Status == UserStatusActive }

type UserProfile struct {
	ID        uuid.UUID  `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	AvatarURL string     `json:"avatar_url,omitempty"`
	Roles     []UserRole `json:"roles"`
	CreatedAt time.Time  `json:"created_at"`
}

func (u *User) ToProfile() UserProfile {
	return UserProfile{ID: u.ID, FirstName: u.FirstName, LastName: u.LastName, AvatarURL: u.AvatarURL, Roles: u.Roles, CreatedAt: u.CreatedAt}
}

type RegisterRequest struct {
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Role      UserRole `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
