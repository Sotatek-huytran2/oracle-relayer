package util

import (
	"testing"

	ethcmm "github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"
)

func TestAlertConfig(t *testing.T) {
	cases := []struct {
		config *AlertConfig
		result bool
	}{
		{
			&AlertConfig{
				Moniker: "",
			},
			true,
		}, {
			&AlertConfig{
				BlockUpdateTimeOut: 0,
			},
			true,
		}, {
			&AlertConfig{
				PackageDelayAlertThreshold: 0,
			},
			true,
		}, {
			&AlertConfig{
				Moniker:                    "test",
				BlockUpdateTimeOut:         10,
				PackageDelayAlertThreshold: 10,
			},
			false,
		},
	}

	for _, config := range cases {
		if config.result {
			require.Panics(t, config.config.Validate, "the check should panic")
		} else {
			require.NotPanics(t, config.config.Validate, "the check should not panic")
		}
	}
}

func TestDBConfig(t *testing.T) {
	cases := []struct {
		config *DBConfig
		result bool
	}{
		{
			&DBConfig{
				Dialect: "wrong",
			},
			true,
		}, {
			&DBConfig{
				Dialect: "mysql",
				DBPath:  "",
			},
			true,
		}, {
			&DBConfig{
				Dialect: "mysql",
				DBPath:  "path",
			},
			false,
		},
	}

	for _, config := range cases {
		if config.result {
			require.Panics(t, config.config.Validate, "the check should panic")
		} else {
			require.NotPanics(t, config.config.Validate, "the check should not panic")
		}
	}
}

func TestChainConfig(t *testing.T) {
	cases := []struct {
		config *ChainConfig
		result bool
	}{
		{
			&ChainConfig{
				ASCStartHeight: -1,
			},
			true,
		}, {
			&ChainConfig{
				ASCStartHeight: 1,
				ASCProviders:   []string{},
			},
			true,
		}, {
			&ChainConfig{
				ASCStartHeight: 1,
				ASCProviders:   []string{"provider"},
				ASCConfirmNum:  0,
			},
			true,
		}, {
			&ChainConfig{
				ASCStartHeight:               1,
				ASCProviders:                 []string{"provider"},
				ASCConfirmNum:                1,
				ASCCrossChainContractAddress: ethcmm.Address{},
			},
			true,
		}, {
			&ChainConfig{
				ASCStartHeight:               1,
				ASCProviders:                 []string{"provider"},
				ASCConfirmNum:                1,
				ASCCrossChainContractAddress: ethcmm.Address{1},
				AFCRpcAddrs:                  []string{"rpc addr"},
			},
			true,
		}, {
			&ChainConfig{
				ASCStartHeight:               1,
				ASCProviders:                 []string{"provider"},
				ASCConfirmNum:                1,
				ASCCrossChainContractAddress: ethcmm.Address{1},
				AFCRpcAddrs:                  []string{"rpc addr"},
				AFCKeyType:                   "wrong",
			},
			true,
		}, {
			&ChainConfig{
				ASCStartHeight:               1,
				ASCProviders:                 []string{"provider"},
				ASCConfirmNum:                1,
				ASCCrossChainContractAddress: ethcmm.Address{1},
				AFCRpcAddrs:                  []string{"rpc addr"},
				AFCKeyType:                   KeyTypeAWSMnemonic,
				AFCAWSRegion:                 "",
			},
			true,
		}, {
			&ChainConfig{
				ASCStartHeight:               1,
				ASCProviders:                 []string{"provider"},
				ASCConfirmNum:                1,
				ASCCrossChainContractAddress: ethcmm.Address{1},
				AFCRpcAddrs:                  []string{"rpc addr"},
				AFCKeyType:                   KeyTypeAWSMnemonic,
				AFCAWSRegion:                 "region",
				AFCAWSSecretName:             "",
			},
			true,
		}, {
			&ChainConfig{
				ASCStartHeight:               1,
				ASCProviders:                 []string{"provider"},
				ASCConfirmNum:                1,
				ASCCrossChainContractAddress: ethcmm.Address{1},
				AFCRpcAddrs:                  []string{"rpc addr"},
				AFCKeyType:                   KeyTypeMnemonic,
				AFCMnemonic:                  "",
			},
			true,
		}, {
			&ChainConfig{
				ASCStartHeight:               1,
				ASCProviders:                 []string{"provider"},
				ASCConfirmNum:                1,
				ASCCrossChainContractAddress: ethcmm.Address{1},
				AFCRpcAddrs:                  []string{"rpc addr"},
				AFCKeyType:                   KeyTypeMnemonic,
				AFCMnemonic:                  "mnemonic",
				RelayInterval:                0,
			},
			true,
		}, {
			&ChainConfig{
				ASCStartHeight:               1,
				ASCProviders:                 []string{"provider"},
				ASCConfirmNum:                1,
				ASCCrossChainContractAddress: ethcmm.Address{1},
				AFCRpcAddrs:                  []string{"rpc addr"},
				AFCKeyType:                   KeyTypeMnemonic,
				AFCMnemonic:                  "mnemonic",
				RelayInterval:                1,
			},
			false,
		},
	}

	for _, config := range cases {
		if config.result {
			require.Panics(t, config.config.Validate, "the check should panic")
		} else {
			require.NotPanics(t, config.config.Validate, "the check should not panic")
		}
	}
}

func TestLogConfig(t *testing.T) {
	cases := []struct {
		config *LogConfig
		result bool
	}{
		{
			&LogConfig{
				UseFileLogger: true,
				Filename:      "",
			},
			true,
		}, {
			&LogConfig{
				UseFileLogger:   true,
				Filename:        "file",
				MaxFileSizeInMB: 0,
			},
			true,
		}, {
			&LogConfig{
				UseFileLogger:        true,
				Filename:             "file",
				MaxFileSizeInMB:      1,
				MaxBackupsOfLogFiles: 0,
			},
			true,
		},
	}

	for _, config := range cases {
		if config.result {
			require.Panics(t, config.config.Validate, "the check should panic")
		} else {
			require.NotPanics(t, config.config.Validate, "the check should not panic")
		}
	}
}
