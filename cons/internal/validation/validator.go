package validation

import (
	"time"

	"l0/cons/internal/models"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("datetime_rfc3339", func(fl validator.FieldLevel) bool {
		_, err := time.Parse(time.RFC3339, fl.Field().String())
		return err == nil
	})
	validate.RegisterStructValidation(deliveryStructLevelValidation, models.Delivery{})
}

func deliveryStructLevelValidation(sl validator.StructLevel) {
	delivery := sl.Current().Interface().(models.Delivery)

	if delivery.Phone == "" && delivery.Email == "" {
		sl.ReportError(delivery.Phone, "Phone", "phone", "phoneOrEmail", "")
		sl.ReportError(delivery.Email, "Email", "email", "phoneOrEmail", "")
	}
}

func Validate(order *models.Order) error {
	return validate.Struct(order)
}
