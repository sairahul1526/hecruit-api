package util

import (
	CONFIG "hecruit-backend/config"
	CONSTANT "hecruit-backend/constant"
	LOGGER "hecruit-backend/logger"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// IsAccessToken - check if token is of access type
// error if refresh type
func IsAccessToken(token string) error {
	data, err := ParseJWTToken(token)
	if err != nil {
		return err
	}
	if data["refresh"] != nil && strings.EqualFold(data["refresh"].(string), "1") {
		return CONSTANT.JWTNotAcsessTokenError
	}
	return nil
}

// GetUserIDFromJWTToken - get user id from token
func GetUserIDFromJWTToken(token string) (string, error) {
	data, err := ParseJWTToken(token)
	if err != nil {
		return "", err
	}
	return data["user_id"].(string), nil
}

// CreateAccessToken - jwt token for accessing api
func CreateAccessToken(data map[string]interface{}) (string, error) {
	return createJWTToken(data, CONSTANT.JWTAccessExpiry, false)
}

// CreateRefreshToken - jwt token for getting access token, if expired
func CreateRefreshToken(data map[string]interface{}) (string, error) {
	return createJWTToken(data, CONSTANT.JWTRefreshExpiry, true)
}

func createJWTToken(data map[string]interface{}, expiry int, refreshToken bool) (string, error) {
	claims := jwt.MapClaims{}
	claims = data
	claims["exp"] = strconv.FormatInt(time.Now().Add(time.Minute*time.Duration(expiry)).Unix(), 10)
	if refreshToken {
		claims["refresh"] = "1"
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString(CONFIG.JWTSecret)
	if err != nil {
		LOGGER.Log("createJWTToken", data, expiry, refreshToken, err)
		return "", err
	}
	return token, nil
}

func ExtractJWTToken(authorization string) string {
	bearToken := authorization
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyJWTToken(authorization string) (*jwt.Token, error) {
	token, err := jwt.Parse(ExtractJWTToken(authorization), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, CONSTANT.JWTUnexpectedSigningMethodError
		}
		return CONFIG.JWTSecret, nil
	})
	if err != nil {
		LOGGER.Log("verifyJWTToken", authorization, err)
		return nil, err
	}
	// check if token valid
	if _, ok := token.Claims.(jwt.Claims); !ok {
		LOGGER.Log("verifyJWTToken", authorization, CONSTANT.JWTInvalidTokenError)
		return nil, CONSTANT.JWTInvalidTokenError
	}
	if !token.Valid {
		LOGGER.Log("verifyJWTToken", authorization, CONSTANT.JWTInvalidTokenError)
		return nil, CONSTANT.JWTInvalidTokenError
	}
	return token, nil
}

// ParseJWTToken - parse jwt token from auth header
func ParseJWTToken(authorization string) (map[string]interface{}, error) {
	token, err := verifyJWTToken(authorization)
	if err != nil {
		LOGGER.Log("ParseJWTToken", authorization, err)
		return map[string]interface{}{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		LOGGER.Log("ParseJWTToken", authorization, CONSTANT.JWTInvalidTokenError)
		return map[string]interface{}{}, CONSTANT.JWTInvalidTokenError
	}

	// extract expiry
	exp, ok := claims["exp"].(string)
	if !ok {
		LOGGER.Log("ParseJWTToken", authorization, CONSTANT.JWTInvalidTokenExpiryError)
		return map[string]interface{}{}, CONSTANT.JWTInvalidTokenExpiryError
	}
	expiry, err := strconv.ParseInt(exp, 10, 64)
	if err != nil {
		LOGGER.Log("ParseJWTToken", authorization, CONSTANT.JWTInvalidTokenExpiryError)
		return map[string]interface{}{}, CONSTANT.JWTInvalidTokenExpiryError
	}

	// check if token expired
	if expiry < time.Now().Unix() { // expired if less than current time
		LOGGER.Log("ParseJWTToken", authorization, CONSTANT.JWTTokenExpiredError)
		return map[string]interface{}{}, CONSTANT.JWTTokenExpiredError
	}

	return claims, nil
}
