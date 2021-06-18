package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"longstory/graph/model"
	"os"
	"reflect"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var SECRET = os.Getenv("SECRET")

const (
	EXPINTEGER                 = 10
	TIMEDIFF           float64 = 5
	ERR_NEED_NEW_TOKEN         = "error need new token"
	ERR_TOKEN_INVALID          = "error token invalid"
)

//createToken return new token string
func CreateToken(user *model.User) (string, error) {

	claims := createClaimMap(user)

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := newToken.SignedString([]byte(SECRET))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

//ParseMapClaims will cast jwt.Token to model.User (map to struct)
func ParseMapClaims(parsedToken *jwt.Token) (*model.User, error) {

	var mapClaims jwt.MapClaims = parsedToken.Claims.(jwt.MapClaims)

	var claims model.User

	dbByte, err := json.Marshal(mapClaims)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dbByte, &claims)
	if err != nil {
		return nil, err
	}
	return &claims, nil
}

//ParseTokenString will parse token string and return jwt.Token. If token is invalid it will return error ERR_INVALID_TOKEN.
//if token is expired it will return with error message ERR_NEED_NEW_TOKEN
func ParseTokenString(tokenstring *string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(*tokenstring, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, errors.New(ERR_TOKEN_INVALID)
	}
	err = checkTokenRenew(parsedToken)
	if err != nil {
		return nil, err
	}
	return parsedToken, nil
}

//checkTokenRenew will substract now with exp value in token and compare to time range, if it fullfills it return error ERR_NEED_NEW_TOKEN
func checkTokenRenew(token *jwt.Token) error {
	now := time.Now().Unix()
	timeSubNow := (*token).Claims.(jwt.MapClaims)["exp"].(float64) - float64(now)

	if timeSubNow <= TIMEDIFF {
		return errors.New(ERR_NEED_NEW_TOKEN)
	}

	return nil
}

//createClaimMap will create new jwt mapclaims from user struct and return it
func createClaimMap(user *model.User) jwt.MapClaims {
	var claims jwt.MapClaims = make(jwt.MapClaims)

	var userval = reflect.ValueOf(*user)
	var usertype = reflect.TypeOf(*user)

	for i := 0; i < userval.NumField(); i++ {
		fieldName := usertype.Field(i).Name
		fieldValue := userval.Field(i).Interface()
		if userval.Field(i).Kind() == reflect.Int64 {
			claims[fieldName] = fieldValue.(int64)
		} else {
			claims[fieldName] = fieldValue
		}
	}

	claims["exp"] = time.Now().Unix() + EXPINTEGER
	//this code is intended to be place after for loop so that new exp override old exp for refresh token

	return claims
}
