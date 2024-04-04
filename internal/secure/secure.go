package secure

import (
	"context"
	"encoding/base64"
	"io"

	"github.com/ivan1993spb/snake-bot/internal/utils"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Secure struct {
	fs    afero.Fs
	clock utils.Clock
}

func New(fs afero.Fs, clock utils.Clock) *Secure {
	return &Secure{
		fs:    fs,
		clock: clock,
	}
}

func (s *Secure) JwtFromFile(ctx context.Context, path string) (*Jwt, error) {
	key, err := readBase64KeyFile(ctx, s.fs, path)
	if err != nil {
		return nil, err
	}

	return NewJwt(key, s.clock), nil
}

func readBase64KeyFile(ctx context.Context, fs afero.Fs, path string) ([]byte, error) {
	log := utils.GetLogger(ctx).WithField("path", path)
	log.Info("reading signing key from file")

	f, err := fs.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open signing key file")
	}
	defer f.Close()

	dec := base64.NewDecoder(base64.StdEncoding, f)
	key, err := io.ReadAll(dec)
	if err != nil {
		return nil, errors.Wrap(err, "read signing key file")
	}

	return key, nil
}
