package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Disbursment struct {
	bun.BaseModel `bun:"table:dibursements"`

	ID                uint       `bun:"id,pk,nullzero"`
	FieldOfficerID    uint       `bun:"field_officer_id"`
	AggreementFileURL string     `bun:"aggreement_file_url"`
	CreatedAt         time.Time  `bun:"created_at"`
	UpdatedAt         *time.Time `bun:"updated_at,nullzero"`
}
