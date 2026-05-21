package auth

import (
	"github.com/NSDN/neonya/apps/server/internal/config"
)

type UserPublicInfo struct {
	UID       string `json:"uid" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string `json:"username" gorm:"unique;not null"`
	Nickname  string `json:"nickname" gorm:"not null"`
	UserGroup string `json:"userGroup" gorm:"not null"`
	Icon      string `json:"icon"`
}

type User struct {
	UID       string `json:"uid" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string `json:"username" gorm:"unique;not null"`
	Password  string `json:"password" gorm:"not null"`
	Salt      string `json:"salt" gorm:"not null"`
	Nickname  string `json:"nickname" gorm:"not null"`
	UserGroup string `json:"userGroup" gorm:"not null"`
	Icon      string `json:"icon"`
}

type RegisterInfo struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	AccessToken string `json:"accessToken"`
}

const (
	DefaultIcon      = "https://i.imgur.com/SH1uR3f.png"
	DefaultUserGroup = "一面の毛玉"
)

var PasswordMaxIndex = config.PASSWORD_MAX_INDEX
var SaltLength = config.SALT_LENGTH
var BcryptCost = config.BCRYPT_COST
var AuthType = config.AUTHENTICATION_TYPE
