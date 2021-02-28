package helper

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
)

type TokenJWT struct {
	User string
	Username string
	Name string
	Role string
	Scope string
	Expired string
	Token string
}

func (r *TokenJWT) ValidateJWT() (jwt.MapClaims, error) {
	token, err := jwt.Parse(r.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid token")
		}
		return []byte("j24g$a@T8#mHN4%"), nil
	})
	if token.Valid {
		claims, _ := token.Claims.(jwt.MapClaims)
		return claims,nil
	}else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return nil, errors.New("Invalid token:Expired")
		}
	}
	return nil, errors.New("Invalid token")
}




func AuthorizeRole(db *gorm.DB,r *http.Request,role string) error {
	token := TokenJWT{}
	decode,_,err := DecodeJWT(&token,r)
	if err != nil {
		return err
	}
	if role != "*" {
		if decode["role"] != role {
			return errors.New("Unauthorized")
		}
	}
	return nil
}

func DecodeJWT(token *TokenJWT,r *http.Request) (jwt.MapClaims,[]string,error){
	reqToken := r.Header.Get("Authorization")
	splitBearer := strings.Split(reqToken, "Bearer")
	if len(splitBearer) != 2 {
		return nil,nil,errors.New("Invalid token")
	}
	token.Token = strings.TrimSpace(splitBearer[1])
	splitToken := strings.Split(token.Token, ".")
	if len(splitToken) != 3 {
		return nil,nil,errors.New("Unauthorized")
	}
	dataToken,err := token.ValidateJWT()
	if err != nil{
		return nil,nil,err
	}
	return dataToken,splitToken, nil
}