package dto

import "demo/model"

type AdminDto struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func ToAdminDto(admin model.Admin) AdminDto {
	return AdminDto{
		ID:    admin.ID,
		Email: admin.Email,
		Name:  admin.Name,
	}
}
