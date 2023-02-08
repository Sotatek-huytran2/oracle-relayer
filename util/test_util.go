package util

import (
	"io/ioutil"

	"github.com/jinzhu/gorm"

	"github.com/Sotatek-huytran2/oracle-relayer/model"
)

var testConfig = `
{
  "db_config": {
    "dialect": "sqlite3",
    "db_path": "user:password@(host:port)/db_name?charset=utf8&parseTime=True&loc=Local"
  },
  "chain_config": {
    "asc_start_height": 1,
    "asc_providers": ["asc_provider"],
    "asc_confirm_num": 2,
    "asc_cross_chain_contract_address": "0x0000000000000000000000000000000000001004",

    "afc_rpc_addrs": ["afc_rpc_addr"],
    "afc_key_type": "mnemonic",
    "afc_aws_region": "",
    "afc_aws_secret_name": "",
    "afc_mnemonic": "",

    "relay_interval": 1000
  },
  "log_config": {
    "level": "INFO",
    "filename": "",
    "max_file_size_in_mb": 0,
    "max_backups_of_log_files": 0,
    "max_age_to_retain_log_files_in_days": 0,
    "use_console_logger": true,
    "use_file_logger": false,
    "compress": false
  },
  "admin_config": {
    "listen_addr": ":8080"
  },
  "alert_config": {
    "moniker": "moniker",
    "telegram_bot_id": "your_bot_id",
    "telegram_chat_id": "your_chat_id",
    "block_update_time_out": 60
  }
}
`

func GetTestConfig() *Config {
	config := ParseConfigFromJson(testConfig)
	return config
}

func PrepareDB(config *Config) (*gorm.DB, error) {
	config.DBConfig.DBPath = "tmp.db"
	tmpDBFile, err := ioutil.TempFile("", config.DBConfig.DBPath)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(config.DBConfig.Dialect, tmpDBFile.Name())
	if err != nil {
		return nil, err
	}
	model.InitTables(db)
	return db, nil
}
