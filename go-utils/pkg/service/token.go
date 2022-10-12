package service

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTService interface {
	Create(time.Duration, string, jwt.MapClaims) (string, error)
	Validate(string) (jwt.MapClaims, error)
}

type jwtService struct {
	// privateKey Private key for the JWT creation
	privateKey []byte

	// publicKey Public key used for JWT validation
	publicKey []byte
}

// NewJWTService Creates a new JWT service from a privateKey and a publicKey. The keys are expected to be RSA256 with PEM
// encoding.
func NewJWTService(keyDir string) (JWTService, error) {

	privateKey, err := ioutil.ReadFile(keyDir + "/key")
	if err != nil {
		return nil, fmt.Errorf("jwt: private key file: %w", err)
	}

	publicKey, err := ioutil.ReadFile(keyDir + "/key.pub")
	if err != nil {
		return nil, fmt.Errorf("jwt: public key file: %w", err)
	}

	return jwtService{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

// Create creates a JWT token string with a given TTL and content for the claim data. If creation is successful, then
// the token string is returned with no error. Otherwise the string will be empty and the error will be set.
func (j jwtService) Create(ttl time.Duration, issuer string, additionalClaims jwt.MapClaims) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}
	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["iss"] = issuer

	for k, v := range additionalClaims {
		if claims[k] != nil {
			return "", fmt.Errorf("claim value already set: %s", k)
		}
		claims[k] = v
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

// Validate Validates a JWT token string. It uses the public key to parse the token. If validation is successful, then
// the claim data is returned and the error is nil. Otherwise, the interface data will be nil and the error will be set.
func (j jwtService) Validate(token string) (jwt.MapClaims, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return nil, fmt.Errorf("validates: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("incorrect signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validation: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validation: invalid")
	}

	return claims, nil
}
