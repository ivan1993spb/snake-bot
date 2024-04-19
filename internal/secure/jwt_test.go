package secure_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-bot/internal/secure"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

func Test_SecureJWT(t *testing.T) {
	const subject = "user"
	key := []byte("secret")

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(utils.NeverClock.Now().Add(time.Minute)),
		Subject:   subject,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(key)
	require.NoError(t, err)

	jwtVerify := secure.NewJwt(key, utils.NeverClock)
	actualSubject, err := jwtVerify.VerifyToken(ss)
	require.NoError(t, err)
	require.Equal(t, subject, actualSubject)
}
