package models

import (
	"time"

	"github.com/uptrace/bun"
)

type UserRole int

const (
	RoleBorower        UserRole = 1
	RoleFieldValidator UserRole = 2
	RoleFieldOfficer   UserRole = 3
	RoleInvestor       UserRole = 4
)

func (s UserRole) String() string {
	switch s {
	case RoleBorower:
		return "borowwer"
	case RoleFieldValidator:
		return "field_validator"
	case RoleFieldOfficer:
		return "field_officer"
	case RoleInvestor:
		return "role_investor"
	default:
		return "unknown"
	}
}

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID        uint      `bun:"id,pk,nullzero"`
	Name      string    `bun:"name"`
	Email     string    `bun:"email"`
	Role      UserRole  `bun:"role"`
	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
}

// func GetUser(role UserRole) *User {
// 	switch role {
// 	case RoleBorower:
// 		return &User{
// 			ID:        1,
// 			Name:      "Borowwer",
// 			Email:     "borowwer@amartha.id",
// 			Role:      RoleBorower,
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		}
// 	case RoleFieldValidator:
// 		return &User{
// 			ID:        2,
// 			Name:      "Field Validator",
// 			Email:     "validator@amartha.id",
// 			Role:      RoleFieldValidator,
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		}
// 	case RoleFieldOfficer:
// 		return &User{
// 			ID:        3,
// 			Name:      "Field Officer",
// 			Email:     "officer@amartha.id",
// 			Role:      RoleFieldOfficer,
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		}
// 	case RoleInvestor:
// 		return &User{
// 			ID:        4,
// 			Name:      "Investor",
// 			Email:     "investor@amartha.id",
// 			Role:      RoleInvestor,
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		}
// 	default:
// 		return nil
// 	}

// }
