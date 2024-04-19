package core

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"

	"github.com/ivan1993spb/snake-bot/internal/config"
	"github.com/ivan1993spb/snake-bot/internal/models"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

type Storage interface {
	Load(ctx context.Context) (map[int]int, error)
	Save(ctx context.Context, state map[int]int) error
	Type() string
}

func NewStorage(fs afero.Fs, cfg config.Storage) Storage {
	if len(cfg.Path) == 0 {
		return &storageMem{}
	}

	return &storageFs{
		fs:   fs,
		path: cfg.Path,
	}
}

type storageFs struct {
	mux  sync.Mutex
	path string
	fs   afero.Fs
}

func (s *storageFs) Load(ctx context.Context) (map[int]int, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	ctx = utils.WithModule(ctx, "storage")
	log := utils.GetLogger(ctx).WithField("path", s.path)

	log.Info("loading state from file")

	f, err := s.fs.Open(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[int]int{}, nil
		}

		return nil, err
	}
	defer f.Close()

	var games *models.Games
	err = yaml.NewDecoder(f).Decode(&games)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return map[int]int{}, nil
		}

		return nil, err
	}

	return games.ToMapState(), nil
}

func (s *storageFs) Save(ctx context.Context, state map[int]int) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	ctx = utils.WithModule(ctx, "storage")
	log := utils.GetLogger(ctx).WithField("path", s.path)

	log.Info("saving state to file")

	f, err := s.fs.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)

	err = enc.Encode(models.NewGames(state))
	if err != nil {
		return err
	}

	enc.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *storageFs) Type() string {
	return "fs"
}

type storageMem struct {
	mux   sync.Mutex
	state map[int]int
}

func (s *storageMem) Load(ctx context.Context) (map[int]int, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.state, nil
}

func (s *storageMem) Save(ctx context.Context, state map[int]int) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.state = state

	return nil
}

func (s *storageMem) Type() string {
	return "memory"
}
