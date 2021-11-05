package validator

import (
	"github.com/DieGopherLT/refaccionaria-backend/internal/helpers"
	"github.com/DieGopherLT/refaccionaria-backend/internal/models"
	"github.com/asaskevich/govalidator"
)

func IsValidProvider(provider models.Provider) (bool, helpers.Response) {
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
