package db

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var SECRET_KEY string = "SECRET222!!"

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Article struct {
	Model
	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}

type User struct {
	Model
	Email      string `json:"email"`
	Password   string `json:"password"`
	Provider   string `json:"provider"`
	ProviderID int    `json:"provider_id"`
}

func (user *User) GeneratePassword() {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		panic(err)
	}
	user.Password = string(bytes)
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

type SignedDetails struct {
	ID    uint
	Email string
	jwt.StandardClaims
}

func (user *User) Token() (string, error) {
	claims := &SignedDetails{
		ID:    user.ID,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(SECRET_KEY))
}

func (user *User) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*SignedDetails); ok && token.Valid {
		if claims.ExpiresAt < time.Now().Local().Unix() {
			return false, errors.New("token expired")
		}
		return true, nil
	} else {
		return false, errors.New("claim invalid")
	}
}
