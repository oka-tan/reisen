package db

import (
	"time"

	"github.com/uptrace/bun"
)

type Report struct {
	bun.BaseModel `bun:"table:report"`

	Board      string    `bun:"board,pk"`
	PostNumber int64     `bun:"post_number,pk"`
	UserIP     string    `bun:"user_ip,type:inet"`
	CreatedAt  time.Time `bun:"created_at"`
	ReportType string    `bun:"report_type"`
	Comment    string    `bun:"comment,nullzero"`
	Handled    bool      `bun:"handled"`
}
