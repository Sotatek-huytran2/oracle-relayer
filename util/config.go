package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	ethcmm "github.com/ethereum/go-ethereum/common"

	"github.com/aximchain/oracle-relayer/common"
)

const (
	KeyTypeMnemonic    = "mnemonic"
	KeyTypeAWSMnemonic = "aws_mnemonic"
)

type Config struct {
	DBConfig    *DBConfig    `json:"db_config"`
	ChainConfig *ChainConfig `json:"chain_config"`
	LogConfig   *LogConfig   `json:"log_config"`
	AlertConfig *AlertConfig `json:"alert_config"`
	AdminConfig *AdminConfig `json:"admin_config"`
}

func (cfg *Config) Validate() {
	cfg.DBConfig.Validate()
	cfg.ChainConfig.Validate()
	cfg.LogConfig.Validate()
	cfg.AlertConfig.Validate()
}

type AlertConfig struct {
	Moniker string `json:"moniker"`

	TelegramBotId  string `json:"telegram_bot_id"`
	TelegramChatId string `json:"telegram_chat_id"`

	PagerDutyAuthToken string `json:"pager_duty_auth_token"`

	BlockUpdateTimeOut         int64 `json:"block_update_time_out"`
	PackageDelayAlertThreshold int64 `json:"package_delay_alert_threshold"`
}

func (cfg *AlertConfig) Validate() {
	if cfg.Moniker == "" {
		panic("moniker should not be empty")
	}

	if cfg.BlockUpdateTimeOut <= 0 {
		panic("block_update_time_out should be larger than 0")
	}

	if cfg.PackageDelayAlertThreshold <= 0 {
		panic("package_delay_alert_threshold should be larger than 0")
	}
}

type DBConfig struct {
	Dialect string `json:"dialect"`
	DBPath  string `json:"db_path"`
}

func (cfg *DBConfig) Validate() {
	if cfg.Dialect != common.DBDialectMysql && cfg.Dialect != common.DBDialectSqlite3 {
		panic(fmt.Sprintf("only %s and %s supported", common.DBDialectMysql, common.DBDialectSqlite3))
	}
	if cfg.DBPath == "" {
		panic("db path should not be empty")
	}
}

type ChainConfig struct {
	ASCStartHeight               int64          `json:"asc_start_height"`
	ASCProviders                 []string       `json:"asc_providers"`
	ASCConfirmNum                int64          `json:"asc_confirm_num"`
	ASCChainId                   uint16         `json:"asc_chain_id"`
	ASCCrossChainContractAddress ethcmm.Address `json:"asc_cross_chain_contract_address"`

	AFCRpcAddrs      []string `json:"afc_rpc_addrs"`
	AFCMnemonic      string   `json:"afc_mnemonic"`
	AFCKeyType       string   `json:"afc_key_type"`
	AFCAWSRegion     string   `json:"afc_aws_region"`
	AFCAWSSecretName string   `json:"afc_aws_secret_name"`

	RelayInterval int64 `json:"relay_interval"`
}

func (cfg *ChainConfig) Validate() {
	if cfg.ASCStartHeight < 0 {
		panic("asc_start_height should not be less than 0")
	}
	if len(cfg.ASCProviders) == 0 {
		panic("asc_providers should not be empty")
	}
	if cfg.ASCConfirmNum <= 0 {
		panic("asc_confirm_num should be larger than 0")
	}

	// replace asc_confirm_num if it is less than DefaultConfirmNum
	if cfg.ASCConfirmNum <= common.DefaultConfirmNum {
		cfg.ASCConfirmNum = common.DefaultConfirmNum
	}

	var emptyAddr ethcmm.Address
	if cfg.ASCCrossChainContractAddress.String() == emptyAddr.String() {
		panic("asc_token_hub_contract_address should not be empty")
	}

	if len(cfg.AFCRpcAddrs) == 0 {
		panic("afc_rpc_addrs should not be empty")
	}
	if cfg.AFCKeyType != KeyTypeMnemonic && cfg.AFCKeyType != KeyTypeAWSMnemonic {
		panic(fmt.Sprintf("afc_key_type of axim chain only supports %s and %s", KeyTypeMnemonic, KeyTypeAWSMnemonic))
	}
	if cfg.AFCKeyType == KeyTypeAWSMnemonic && cfg.AFCAWSRegion == "" {
		panic("afc_aws_region of axim chain should not be empty")
	}
	if cfg.AFCKeyType == KeyTypeAWSMnemonic && cfg.AFCAWSSecretName == "" {
		panic("afc_aws_secret_name of axim chain should not be empty")
	}
	if cfg.AFCKeyType == KeyTypeMnemonic && cfg.AFCMnemonic == "" {
		panic("afc_mnemonic should not be empty")
	}

	if cfg.RelayInterval <= 0 {
		panic(fmt.Sprintf("relay interval should be larger than 0"))
	}
}

type LogConfig struct {
	Level                        string `json:"level"`
	Filename                     string `json:"filename"`
	MaxFileSizeInMB              int    `json:"max_file_size_in_mb"`
	MaxBackupsOfLogFiles         int    `json:"max_backups_of_log_files"`
	MaxAgeToRetainLogFilesInDays int    `json:"max_age_to_retain_log_files_in_days"`
	UseConsoleLogger             bool   `json:"use_console_logger"`
	UseFileLogger                bool   `json:"use_file_logger"`
	Compress                     bool   `json:"compress"`
}

func (cfg *LogConfig) Validate() {
	if cfg.UseFileLogger {
		if cfg.Filename == "" {
			panic("filename should not be empty if use file logger")
		}
		if cfg.MaxFileSizeInMB <= 0 {
			panic("max_file_size_in_mb should be larger than 0 if use file logger")
		}
		if cfg.MaxBackupsOfLogFiles <= 0 {
			panic("max_backups_off_log_files should be larger than 0 if use file logger")
		}
	}
}

type AdminConfig struct {
	ListenAddr string `json:"listen_addr"`
}

// ParseConfigFromFile returns the config from json file
func ParseConfigFromFile(filePath string) *Config {
	bz, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	var config Config
	if err := json.Unmarshal(bz, &config); err != nil {
		panic(err)
	}
	return &config
}

// ParseConfigFromJson returns the config from json string
func ParseConfigFromJson(content string) *Config {
	var config Config
	if err := json.Unmarshal([]byte(content), &config); err != nil {
		panic(err)
	}
	return &config
}
