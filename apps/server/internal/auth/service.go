package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/NSDN/neonya/apps/server/internal/shared"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(db *gorm.DB, info *RegisterInfo) (*UserPublicInfo, error) {
	_, err := findUserByUsername(db, info.Username)

	if err == nil {
		return nil, errors.New(shared.Messages.AuthorizeFailedUserExist)
	}

	if err.Error() != shared.Messages.AuthorizeFailedNoUser {
		return nil, err
	}

	return createNewUser(db, info)
}

func Login(db *gorm.DB, tokenKey string, info LoginInfo) (string, error) {
	user, err := findUserByUsername(db, info.Username)

	if err != nil {
		return "", err
	}

	err = comparePassword(info.Password, user.Password, user.Salt)

	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.UID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(tokenKey))
}

func FindUserByUID(db *gorm.DB, uid string) (*UserPublicInfo, error) {
	var user User
	result := db.Where("uid = ?", uid).Take(&user)

	if result.Error != nil {
		return nil, errors.New(shared.Messages.AuthorizeFailedNoUser)
	}

	publicInfo := UserPublicInfo{
		UID:       user.UID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		UserGroup: user.UserGroup,
		Icon:      user.Icon,
	}

	return &publicInfo, nil
}

func findUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	result := db.Where("username = ?", username).Take(&user)

	if result.Error != nil {
		return nil, errors.New(shared.Messages.AuthorizeFailedNoUser)
	}

	return &user, nil
}

func createNewUser(db *gorm.DB, info *RegisterInfo) (*UserPublicInfo, error) {
	salt, err := generateBase64Salt()

	if err != nil {
		return nil, err
	}

	encryptedPassword, err := encryptPassword(info.Password, salt)

	if err != nil {
		return nil, err
	}

	user := User{
		Username:  info.Username,
		Password:  encryptedPassword,
		Salt:      salt,
		Nickname:  info.Username,
		Icon:      DefaultIcon,
		UserGroup: DefaultUserGroup,
	}

	result := db.Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	publicInfo := UserPublicInfo{
		UID:       user.UID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		UserGroup: user.UserGroup,
		Icon:      user.Icon,
	}

	return &publicInfo, nil
}

func encryptPassword(password string, salt string) (string, error) {
	byteSalt, err := base64.StdEncoding.DecodeString(salt)

	if err != nil {
		return "", err
	}

	appended := append(
		[]byte(password[:PasswordMaxIndex]),
		byteSalt...,
	)

	hashed, err := bcrypt.GenerateFromPassword(appended, BcryptCost)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hashed), nil
}

func comparePassword(password string, base64Hashed string, base64Salt string) error {
	hashed, err := base64.StdEncoding.DecodeString(base64Hashed)

	if err != nil {
		return err
	}

	salt, err := base64.StdEncoding.DecodeString(base64Salt)

	if err != nil {
		return err
	}

	payload := append(
		[]byte(password[:PasswordMaxIndex]),
		salt...,
	)

	err = bcrypt.CompareHashAndPassword(hashed, payload)

	if err != nil {
		log.Println(err)
		err = errors.New("密码不一致")
	}

	return err
}

func generateBase64Salt() (string, error) {
	salt := make([]byte, SaltLength)

	_, err := rand.Read(salt)

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(salt), nil
}
