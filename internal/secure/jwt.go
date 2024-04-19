package secure

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"

	"github.com/ivan1993spb/snake-bot/internal/utils"
)

type Jwt struct {
	key        []byte
	clock      utils.Clock
	validators map[string]*jwt.Validator
}

func NewJwt(key []byte, clock utils.Clock) *Jwt {
	j := &Jwt{
		key:   key,
		clock: clock,
	}

	j.initValidators()

	return j
}

func (j *Jwt) initValidators() {
	j.validators = map[string]*jwt.Validator{
		"admin": jwt.NewValidator(
			jwt.WithSubject("admin"),
			jwt.WithTimeFunc(j.clock.Now),
		),
		"service": jwt.NewValidator(
			jwt.WithSubject("service"),
			jwt.WithTimeFunc(j.clock.Now),
		),
		"user": jwt.NewValidator(
			jwt.WithSubject("user"),
			jwt.WithExpirationRequired(),
			jwt.WithTimeFunc(j.clock.Now),
		),
	}
}

func (j *Jwt) VerifyToken(tokenString string) (string, error) {
	token, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return "", errors.Wrap(err, "failed to get subject from token")
	}

	validator, ok := j.validators[subject]
	if !ok {
		return "", errors.Errorf("unknown subject %q", subject)
	}

	err = validator.Validate(token.Claims)
	if err != nil {
		return "", errors.Wrap(err, "error validating token")
	}

	return subject, nil
}

func (j *Jwt) ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		j.keyFunc,
		jwt.WithTimeFunc(j.clock.Now),
	)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (j *Jwt) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	if len(j.key) == 0 {
		return nil, errors.New("empty signing key")
	}

	return j.key, nil
}
