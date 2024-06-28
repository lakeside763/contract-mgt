package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Model struct {
	CreatedAt 	time.Time 	`json:"createdAt"`
	UpdateAt 		time.Time 	`json:"updatedAt"`
}

type BaseType string

func (b *BaseType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
		case []byte:
			*b = BaseType(v)
		case string:
			*b = BaseType(v)
		default:
			return fmt.Errorf("cannot scan %T into BaseType", value)
	}
	return nil
}

func (b BaseType) Value() (driver.Value, error) {
	return string(b), nil
}