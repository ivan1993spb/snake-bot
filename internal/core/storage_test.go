package core_test

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-bot/internal/config"
	"github.com/ivan1993spb/snake-bot/internal/core"
)

func Test_Storage(t *testing.T) {
	ctx := context.Background()

	type args struct {
		fs  afero.Fs
		cfg config.Storage
	}

	tests := []struct {
		name       string
		args       args
		expectType string
	}{
		{
			name: "empty path",
			args: args{
				fs:  afero.NewMemMapFs(),
				cfg: config.Storage{},
			},
			expectType: "memory",
		},
		{
			name: "non-empty path",
			args: args{
				fs: afero.NewMemMapFs(),
				cfg: config.Storage{
					Path: "/test",
				},
			},
			expectType: "fs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := core.NewStorage(tt.args.fs, tt.args.cfg)
			require.Equal(t, tt.expectType, storage.Type())

			t.Run("Save", func(t *testing.T) {
				err := storage.Save(ctx, map[int]int{
					1: 2,
					2: 3,
					3: 4,
				})
				require.NoError(t, err)
			})

			t.Run("Load", func(t *testing.T) {
				state, err := storage.Load(ctx)
				require.NoError(t, err)
				require.Equal(t, map[int]int{
					1: 2,
					2: 3,
					3: 4,
				}, state)
			})
		})
	}
}

func Test_storageFs_Load_NoFile(t *testing.T) {
	ctx := context.Background()
	fs := afero.NewMemMapFs()

	storage := core.NewStorage(fs, config.Storage{
		Path: "/test/some_config_file",
	})

	state, err := storage.Load(ctx)
	require.NoError(t, err)
	require.Empty(t, state)
}

func Test_storageFs_Load_EmptyFile(t *testing.T) {
	const filePath = "/var/lib/snake-bot/some_config_file"

	ctx := context.Background()
	fs := afero.NewMemMapFs()

	require.NoError(t, fs.MkdirAll("/var/lib/snake-bot", 0700))
	require.NoError(t, afero.WriteFile(fs, filePath, []byte{}, 0600))

	storage := core.NewStorage(fs, config.Storage{
		Path: filePath,
	})

	state, err := storage.Load(ctx)
	require.NoError(t, err)
	require.Empty(t, state)
}

func Test_storageFs_Save_FileExists(t *testing.T) {
	const filePath = "/test/some_config_file"

	ctx := context.Background()
	fs := afero.NewMemMapFs()

	storage := core.NewStorage(fs, config.Storage{
		Path: filePath,
	})

	err := storage.Save(ctx, map[int]int{
		1: 2,
		2: 3,
		3: 4,
	})
	require.NoError(t, err)

	_, err = fs.Stat(filePath)
	require.NoError(t, err)
}
