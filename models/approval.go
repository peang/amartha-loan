package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Approval struct {
	bun.BaseModel `bun:"table:approvals"`

	ID               uint       `bun:"id,pk,nullzero"`
	FieldValidatorID uint       `bun:"field_validator_id"`
	ApprovalFileURL  string     `bun:"approval_file_url"`
	CreatedAt        time.Time  `bun:"created_at"`
	UpdatedAt        *time.Time `bun:"updated_at,nullzero"`
}
