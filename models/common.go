package models

import (
	"database/sql/driver"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var JwtKey []byte

// Initialize validator and translator
var (
	Validate *validator.Validate
	Trans    ut.Translator
)

func Init() {
	jwtSecret := os.Getenv("JWT_SECRET")
	JwtKey = []byte(jwtSecret)

	Validate = validator.New()

	// Initialize the translator for English locale
	eng := en.New()
	uni := ut.New(eng, eng)
	Trans, _ = uni.GetTranslator("en")

	// Register translations for English
	enTranslations.RegisterDefaultTranslations(Validate, Trans)
}

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