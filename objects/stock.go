package objects

import (
	"time"
)

type Stock struct {
	// Identifier
	ID string `gorm:"primary_key" json:"id,omitempty"`

	// General details
	Name  string  `json:"name,omitempty"`
	Price float64 `json:"price,omitempty"`

	Availability int       `json:"availability,omitempty"`
	IsActive     bool      `json:"is_active,omitempty"`
	CreatedOn    time.Time `json:"created_on,omitempty"`
	UpdatedOn    time.Time `json:"updated_on,omitempty"`
}
