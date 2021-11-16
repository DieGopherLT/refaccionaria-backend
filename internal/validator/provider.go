package validator

import (
	"github.com/DieGopherLT/refaccionaria-backend/internal/helpers"
	"github.com/DieGopherLT/refaccionaria-backend/internal/models"
	"github.com/asaskevich/govalidator"
)

// IsValidProvider checks if a incoming provider has it's phone and email in a valid format
func IsValidProvider(provider models.ProviderDTO) (bool, helpers.Response) {
	isEmail := govalidator.IsEmail(provider.Email)
	if !isEmail {
		resp := helpers.Response{Message: "Correo no válido", Error: true}
		return false, resp
	}

	isPhone := govalidator.IsNumeric(provider.Phone)
	if !isPhone {
		resp := helpers.Response{Message: "Teléfono no válido", Error: true}
		return false, resp
	}

	return true, helpers.Response{}
}
