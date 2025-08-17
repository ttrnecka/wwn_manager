package mapper

import (
	"github.com/ttrnecka/wwn_identity/webapi/internal/entity"
	"github.com/ttrnecka/wwn_identity/webapi/shared/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToUserEntity(u dto.UserDTO) entity.User {
	user := entity.User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
	user.ID, _ = primitive.ObjectIDFromHex(u.ID)
	return user
}

func ToUserDTO(user entity.User) dto.UserDTO {
	return dto.UserDTO{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
	}
}
