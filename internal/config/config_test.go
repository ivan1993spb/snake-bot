package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DefaultConfig_ReturnsDefaultConfig(t *testing.T) {
	require.Equal(t, defaultConfig, DefaultConfig())
}

func Test_ParseFlags_ParsesFlagsCorrectly(t *testing.T) {
	const flagSetName = "test"

	type Test struct {
		msg string

		args     []string
		defaults Config

		expectConfig Config
		expectErr    bool
	}

	var tests = make([]*Test, 0)

	// Test case 1
	tests = append(tests, &Test{
		msg: "run without arguments",

		args:     []string{},
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    false,
	})

	// Test case 2
	configTest2 := defaultConfig
	configTest2.Server.Address = ":7070"

	tests = append(tests, &Test{
		msg: "change address",

		args: []string{
			"-address", ":7070",
		},
		defaults: defaultConfig,

		expectConfig: configTest2,
		expectErr:    false,
	})

	// Test case 3
	configTest3 := defaultConfig
	configTest3.Server.Address = "localhost:6670"
	configTest3.Target.Address = "any-snake-dot.com"

	tests = append(tests, &Test{
		msg: "change address and server address",

		args: []string{
			"-address", "localhost:6670",
			"-snake-server", "any-snake-dot.com",
		},
		defaults: defaultConfig,

		expectConfig: configTest3,
		expectErr:    false,
	})

	// Test case 4
	configTest4 := defaultConfig
	configTest4.Server.Address = "snakeonline.xyz:7986"
	configTest4.Log.EnableJSON = true
	configTest4.Target.WSS = true

	tests = append(tests, &Test{
		msg: "change address, logging and protocol",

		args: []string{
			"-address", "snakeonline.xyz:7986",
			"-log-json",
			"-wss",
		},
		defaults: defaultConfig,

		expectConfig: configTest4,
		expectErr:    false,
	})

	// Test case 5
	configTest5 := defaultConfig
	configTest5.Server.Address = "snakeonline.xyz:3211"
	configTest5.Log.EnableJSON = true
	configTest5.Bots.Limit = 541
	configTest5.Server.Debug = true

	tests = append(tests, &Test{
		msg: "change address, logging, bots limit and debug",

		args: []string{
			"-address", "snakeonline.xyz:3211",
			"-log-json",
			"-bots-limit", "541",
			"-debug",
		},
		defaults: defaultConfig,

		expectConfig: configTest5,
		expectErr:    false,
	})

	// Test case 6
	tests = append(tests, &Test{
		msg: "change address, and make 1 mistake",

		args: []string{
			"-address", "snakeonline.xyz:3211",
			"-foobar", // unknown flag
		},
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    true,
	})

	// Test case 7
	tests = append(tests, &Test{
		msg: "change address, enable debug, make 2 mistakes",

		args: []string{
			"-address", "snakeonline.xyz:3211",
			"-bots-limit", "error", // should be a number
			"-foobar",
		},
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    true,
	})

	// Test case 8
	tests = append(tests, &Test{
		msg: "args is nil",

		args:     nil,
		defaults: defaultConfig,

		expectConfig: defaultConfig,
		expectErr:    false,
	})

	for n, test := range tests {
		t.Log(test.msg)

		label := fmt.Sprintf("case number %d", n+1)
		flagSet := flag.NewFlagSet(flagSetName, flag.ContinueOnError)
		flagSet.SetOutput(ioutil.Discard)

		config, err := ParseFlags(flagSet, test.args, test.defaults)

		if test.expectErr {
			require.NotNil(t, err, label)
		} else {
			require.Nil(t, err, label)
		}
		require.Equal(t, test.expectConfig, config, label)
	}
}

func Test_Config_Fields_ReturnsFieldsOfTheConfig(t *testing.T) {
	require.Equal(t, map[string]interface{}{
		fieldLabelAddress:    ":9999",
		fieldLabelForbidCORS: true,
		fieldLabelDebug:      true,

		fieldLabelSnakeServer: "localhost:9210",
		fieldLabelWSS:         false,

		fieldLabelBotsLimit: 1337,

		fieldLabelLogEnableJSON: false,
		fieldLabelLogLevel:      "warning",
	}, Config{
		Server: Server{
			Address: ":9999",

			ForbidCORS: true,
			Debug:      true,
		},

		Target: Target{
			Address: "localhost:9210",
			WSS:         false,
		},

		Bots: Bots{
			Limit: 1337,
		},

		Log: Log{
			EnableJSON: false,
			Level:      "warning",
		},
	}.Fields())
}
