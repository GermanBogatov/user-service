package entity

import (
	"github.com/GermanBogatov/user-service/internal/config"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          string
	Name        string
	Surname     string
	Email       string
	Password    string
	Roles       []string
	CreatedDate time.Time
	UpdatedDate *time.Time
	JWT         *JWT
}

func (u *User) GenerateID() {
	u.ID = uuid.New().String()
}

func (u *User) GenerateCreatedDate() {
	u.CreatedDate = time.Now().UTC()
}

func (u *User) AddRoleDeveloper() {
	u.Roles = append(u.Roles, config.RoleDeveloper)
}

func (u *User) AddRoleAdmin() {
	u.Roles = append(u.Roles, config.RoleAdmin)
}

func (u *User) SetPasswordHash(hash string) {
	u.Password = hash
}

func (u *User) SetJWT(token, refreshToken string) {
	u.JWT = &JWT{
		Token:        token,
		RefreshToken: refreshToken,
	}
}

type JWT struct {
	Token        string
	RefreshToken string
}
