package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string , error){
	hash , err := bcrypt.GenerateFromPassword([]byte(password) , 10)
	if err != nil {
		return "" , err
	}
	return string(hash) , nil
}

func CheckPasswordHash(hashedPassword , password string) error{
		return bcrypt.CompareHashAndPassword([]byte(hashedPassword) , []byte(password))
}
func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error){
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy", 
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		Subject: userID.String(),
	}	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256 , claims)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
	parsedToken , err := jwt.ParseWithClaims(tokenString , & jwt.RegisteredClaims{} , func(token *jwt.Token)(interface{} , error){
		return []byte(tokenSecret) , nil
	})
	if err != nil {
		return uuid.UUID{} , err
	}
	claims , ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok || !parsedToken.Valid {
		return uuid.UUID{}, errors.New("invalid token")
	}

	uuidStr , err := claims.GetSubject();
	id , err:= uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{} , err
	}
	return id , nil 
}

func GetBearerToken(headers http.Header) (string, error){
	authHeader := headers.Get("Authorization")

	if authHeader == ""{
		return "" , errors.New("Authorization Header does not exist")
	}
	parts := strings.Split(authHeader , " ")
	
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer"{
		return "" , errors.New("Authorization Header does not exist")
	}
	return parts[1] , nil
}

func MakeRefreshToken() (string, error){
	randNum := make([]byte , 32)
	_ , err := rand.Read(randNum)
	if err != nil{
		return "" , err
	}
	return hex.EncodeToString(randNum) , nil
}

func GetAPIKey(headers http.Header) (string, error){
	authHeader := headers.Get("Authorization")

	if authHeader == ""{
		return "" , errors.New("Authorization Header does not exist")
	}
	parts := strings.Split(authHeader , " ")
	
	if len(parts) != 2 || strings.ToLower(parts[0]) != "apikey"{
		return "" , errors.New("Authorization Header does not exist")
	}
	return parts[1] , nil
}