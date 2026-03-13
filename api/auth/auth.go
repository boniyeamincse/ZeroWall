package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	IssuedAt int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	ID       string `json:"jti"`
}

type JWT struct {
	secret []byte
}

func NewJWT(secret string) *JWT {
	return &JWT{secret: []byte(secret)}
}

func (j *JWT) generateSignature(header string, payload string) string {
	message := header + "." + payload
	h := hmac.New(sha256.New, j.secret)
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func (j *JWT) verifySignature(message string, signature string) bool {
	expectedSig := j.generateSignature(strings.Split(message, ".")[0], strings.Split(message, ".")[1])
	return subtle.ConstantTimeCompare([]byte(expectedSig), []byte(signature)) == 1
}

func (j *JWT) GenerateToken(username string, role string, ttl time.Duration) (string, error) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	now := time.Now()
	claims := Claims{
		Username: username,
		Role:     role,
		IssuedAt: now.Unix(),
		ExpiresAt: now.Add(ttl).Unix(),
		ID:       fmt.Sprintf("%d-%s", now.UnixNano(), username),
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signature := j.generateSignature(header, payload)

	return header + "." + payload + "." + signature, nil
}

func (j *JWT) VerifyToken(tokenString string) (*Claims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	header, payload, signature := parts[0], parts[1], parts[2]
	message := header + "." + payload

	if !j.verifySignature(message, signature) {
		return nil, ErrInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}

func (j *JWT) RefreshToken(tokenString string, ttl time.Duration) (string, error) {
	claims, err := j.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}
	return j.GenerateToken(claims.Username, claims.Role, ttl)
}

func (j *JWT) ExtractTokenFromHeader(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}
	return parts[1], nil
}

func ValidatePassword(password, hash string) bool {
	passwordHash := HashPassword(password)
	return subtle.ConstantTimeCompare([]byte(passwordHash), []byte(hash)) == 1
}

func HashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(h[:])
}

type User struct {
	Username string
	Password string
	Role     string
}

func AuthenticateUser(username, password string, users []User) (*User, error) {
	for _, user := range users {
		if user.Username == username && ValidatePassword(password, user.Password) {
			return &user, nil
		}
	}
	return nil, errors.New("invalid credentials")
}
