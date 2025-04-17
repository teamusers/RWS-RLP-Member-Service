package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// Date is a custom type based on time.Time that works with JSON and SQL.
type Date time.Time

// UnmarshalJSON parses a JSON string in the "2006-01-02" format,
// while gracefully handling empty strings and JSON null.
func (d *Date) UnmarshalJSON(data []byte) error {
	// Check if the JSON value is null.
	if string(data) == "null" {
		*d = Date(time.Time{})
		return nil
	}

	// Remove the surrounding quotes.
	s := strings.Trim(string(data), "\"")
	// Return zero value if the string is empty.
	if s == "" {
		*d = Date(time.Time{})
		return nil
	}

	// Parse the non-empty date string.
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// MarshalJSON outputs the date in "2006-01-02" format.
// If the date is the zero value, it outputs JSON null.
func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return []byte("null"), nil
	}
	formatted := t.Format("2006-01-02")
	return []byte(fmt.Sprintf("\"%s\"", formatted)), nil
}

// Value implements the driver.Valuer interface.
// It returns nil (SQL NULL) if the Date is the zero value.
func (d Date) Value() (driver.Value, error) {
	t := time.Time(d)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

// Scan implements the sql.Scanner interface.
// It converts a database value into a Date, handling NULL values.
func (d *Date) Scan(value interface{}) error {
	// If the DB column is NULL, assign zero value.
	if value == nil {
		*d = Date(time.Time{})
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan Date: %v", value)
	}
	*d = Date(t)
	return nil
}

// User represents a user model that maps to a MySQL table.
type User struct {
	ID                    uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ExternalID            string    `gorm:"column:external_id" json:"external_id"`
	ExternalTYPE          string    `gorm:"column:external_id_type" json:"external_id_type"`
	Email                 string    `gorm:"column:email;size:255" json:"email"`
	BurnPin               uint64    `gorm:"column:burn_pin" json:"burn_pin"`
	SessionToken          string    `gorm:"column:session_token" json:"session_token"`
	SessionExpiry         int64     `gorm:"column:session_expiry" json:"session_expiry"`
	GR_ID                 string    `gorm:"column:gr_id" json:"gr_id"`
	RLP_ID                string    `gorm:"column:rlp_id" json:"rlp_id"`
	RWS_Membership_ID     string    `gorm:"column:rws_membership_id" json:"rws_membership_id"`
	RWS_Membership_Number uint64    `gorm:"column:rws_membership_number" json:"rws_membership_number"`
	CreatedAt             time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the table name for the User model.
func (User) TableName() string {
	return "users"
}
