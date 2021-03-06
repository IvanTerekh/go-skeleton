package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/ivanterekh/go-skeleton/internal/env"
	"github.com/ivanterekh/go-skeleton/internal/users"
)

const exp = time.Hour * 48

var method = jwt.SigningMethodHS256

// Authenticator is a service for logging in and
// authenticating users.
type Authenticator struct {
	exp    time.Duration
	method *jwt.SigningMethodHMAC
	secret string
	users  users.Repository
}

// Exp returns expiry time.
func (a *Authenticator) Exp() time.Duration {
	return a.exp
}

// NewAuthenticator returns a new authenticator instance.
func NewAuthenticator(expiry time.Duration, signingMethod *jwt.SigningMethodHMAC,
	secret string, users users.Repository) *Authenticator {

	return &Authenticator{
		exp:    expiry,
		method: signingMethod,
		secret: secret,
		users:  users,
	}
}

// DefaultAuthenticator returns a new authenticator
// initialized with default and global parameters.
func DefaultAuthenticator() *Authenticator {
	return NewAuthenticator(
		exp,
		method,
		env.GetString("AUTH_SECRET", "secret"),
		users.NewMockRepository(),
	)
}

// SetUsersRepo is setter for users repository.
func (a *Authenticator) SetUsersRepo(repo users.Repository) {
	a.users = repo
}

// GenToken generates a new token if user with
// provided credentials exists.
func (a *Authenticator) GenToken(email, password string) (string, error) {
	user, err := a.users.GetByCreds(email, password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(a.method, jwt.MapClaims{
		"exp": time.Now().Add(a.exp),
		"sub": user.ID,
	})

	return token.SignedString([]byte(a.secret))
}

// Authenticate returns user if token is valid.
func (a *Authenticator) Authenticate(tokenStr string) (*users.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not parse token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("could not convert claims to jwt.MapClaims")
	}

	userID, err := getUserID(claims)
	if err != nil {
		return nil, errors.Wrap(err, "could not get user id from claims")
	}

	return a.users.GetByID(userID)
}

func getUserID(claims jwt.MapClaims) (int, error) {
	userIDValue, ok := claims["sub"]
	if !ok {
		return 0, errors.New("claims do not contain user id")
	}

	userID, ok := userIDValue.(float64)
	if !ok {
		return 0, errors.New("user id is not a number")
	}
	return int(userID), nil
}
