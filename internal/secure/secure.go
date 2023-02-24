package secure

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"sync"

	"github.com/pkg/errors"
)

type Secure struct {
	mux *sync.Mutex
	sum [sha256.Size]byte
}

func NewSecure() *Secure {
	return &Secure{
		mux: &sync.Mutex{},
	}
}

const bufferSize = 48

var (
	messageTokenBegin = []byte("Auth token: ")
	messageTokenEnd   = []byte("\n")
)

func (s *Secure) GenerateToken(w io.Writer) error {
	buffer := make([]byte, bufferSize)

	// TODO: add noescape operators.
	if _, err := rand.Reader.Read(buffer); err != nil {
		return errors.Wrap(err, "generate token")
	}

	sum := sha256.Sum256(buffer)
	s.mux.Lock()
	s.sum = sum
	s.mux.Unlock()

	if _, err := w.Write(messageTokenBegin); err != nil {
		return errors.Wrap(err, "fail to warn")
	}
	enc := base64.NewEncoder(base64.RawURLEncoding, w)
	if _, err := enc.Write(buffer); err != nil {
		enc.Close()
		return errors.Wrap(err, "write auth token")
	}
	enc.Close()
	if _, err := w.Write(messageTokenEnd); err != nil {
		return errors.Wrap(err, "fail to warn")
	}

	return nil
}

func (s *Secure) VerifyToken(token string) bool {
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	flag := false
	sum := sha256.Sum256(decoded)
	s.mux.Lock()
	flag = sum == s.sum
	s.mux.Unlock()

	return flag
}
