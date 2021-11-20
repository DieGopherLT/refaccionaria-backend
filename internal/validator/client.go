package validator

import (
	"github.com/DieGopherLT/refaccionaria-backend/internal/helpers"
	"github.com/DieGopherLT/refaccionaria-backend/internal/models"
	"github.com/asaskevich/govalidator"
)

func IsValidClient(client models.ClientDTO) (bool, helpers.Response) {
	isPhone := govalidator.IsNumeric(client.Phone)
	if !isPhone {
		resp := helpers.Response{Message: "Teléfono no válido", Error: true}
		return false, resp
	}

	return true, helpers.Response{}
}
