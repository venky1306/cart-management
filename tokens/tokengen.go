package tokens

import (
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY string = os.Getenv("SECRET_KEY")

type Claims struct {
	Email      string
	First_Name string
	Last_Name  string
	User_Type  string
	User_ID    string
	jwt.RegisteredClaims
}

func GenerateAllTokens(email, firstname, lastname, usertype, userid string) (string, string, error) {

	claims := Claims{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		User_Type:  usertype,
		User_ID:    userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	refreshclaims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 30 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	return token, refreshToken, nil
}

func UpdateTokens() error {

}

func ValidateToken(tokenString string) (claims *Claims, msg string) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) { return []byte(SECRET_KEY), nil })
	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		msg = fmt.Sprintf("token invalid")
		return
	}

	if (*claims.ExpiresAt).Before(time.Now()) {
		msg = fmt.Sprint("token is expired")
		return
	}
	return claims, msg
}
