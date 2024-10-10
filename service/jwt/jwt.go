package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWTService interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	ValidateToken(token string, isAccessToken bool) (*jwt.Token, error)
	GetTokenExpired(token string, isAccessToken bool) (time.Time, error)
}

type JWTCustomClaim struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtService struct {
	accessSecret  string
	refreshSecret string
	accessExp     int
	refreshExp    int
	issuer        string
}

func NewJWTService(cfg *viper.Viper) JWTService {
	return &jwtService{
		accessSecret:  cfg.GetString("JWT_ACCESS_SECRET"),
		refreshSecret: cfg.GetString("JWT_REFRESH_SECRET"),
		accessExp:     cfg.GetInt("JWT_ACCESS_EXP"),
		refreshExp:    cfg.GetInt("JWT_REFRESH_EXP"),
		issuer:        cfg.GetString("JWT_ISSUER"),
	}
}

func (s *jwtService) GenerateAccessToken(userID string) (string, error) {
	claims := &JWTCustomClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(s.accessExp))),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.accessSecret))
}

func (s *jwtService) GenerateRefreshToken(userID string) (string, error) {
	claims := &JWTCustomClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(s.refreshExp))),
			Issuer:    s.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshSecret))
}

func (s *jwtService) ValidateToken(tokenString string, isAccessToken bool) (*jwt.Token, error) {
	secret := s.accessSecret
	if !isAccessToken {
		secret = s.refreshSecret
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if _, ok := token.Claims.(*JWTCustomClaim); ok && token.Valid {
		return token, nil
	} else {
		return nil, err
	}
}

// Function to get when the token will expire
func (s *jwtService) GetTokenExpired(tokenString string, isAccessToken bool) (time.Time, error) {
	token, err := s.ValidateToken(tokenString, isAccessToken)
	if err != nil {
		return time.Time{}, err
	}

	claims, ok := token.Claims.(*JWTCustomClaim)
	if !ok {
		return time.Time{}, err
	}

	return claims.ExpiresAt.Time, nil
}
