package main

import (
	"flag"
	"fmt"

	"github.com/aximchain/go-sdk/common/types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/Sotatek-huytran2/oracle-relayer/admin"
	"github.com/Sotatek-huytran2/oracle-relayer/executor/afc"
	"github.com/Sotatek-huytran2/oracle-relayer/executor/asc"
	"github.com/Sotatek-huytran2/oracle-relayer/model"
	"github.com/Sotatek-huytran2/oracle-relayer/observer"
	"github.com/Sotatek-huytran2/oracle-relayer/relayer"
	"github.com/Sotatek-huytran2/oracle-relayer/util"
)

const (
	flagConfigType         = "config-type"
	flagConfigAwsRegion    = "aws-region"
	flagConfigAwsSecretKey = "aws-secret-key"
	flagConfigPath         = "config-path"
	flagAFCNetwork         = "afc-network"
)

const (
	ConfigTypeLocal = "local"
	ConfigTypeAws   = "aws"
)

func initFlags() {
	flag.String(flagConfigPath, "", "config path")
	flag.String(flagConfigType, "", "config type, local or aws")
	flag.String(flagConfigAwsRegion, "", "aws s3 region")
	flag.String(flagConfigAwsSecretKey, "", "aws s3 secret key")
	flag.Int(flagAFCNetwork, int(types.TestNetwork), "afc chain network type")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		panic(fmt.Sprintf("bind flags error, err=%s", err))
	}
}

func printUsage() {
	fmt.Print("usage: ./relayer --afc-network [0 for testnet, 1 for mainnet] --config-type [local or aws] --config-path config_file_path\n")
}

func main() {
	initFlags()

	afcNetwork := viper.GetInt(flagAFCNetwork)
	if afcNetwork != int(types.TestNetwork) &&
		afcNetwork != int(types.ProdNetwork) &&
		afcNetwork != int(types.TmpTestNetwork) &&
		afcNetwork != int(types.GangesNetwork) {
		printUsage()
		return
	}

	types.Network = types.ChainNetwork(afcNetwork)

	configType := viper.GetString(flagConfigType)
	if configType == "" {
		printUsage()
		return
	}

	if configType != ConfigTypeAws && configType != ConfigTypeLocal {
		printUsage()
		return
	}

	var config *util.Config
	if configType == ConfigTypeAws {
		awsSecretKey := viper.GetString(flagConfigAwsSecretKey)
		if awsSecretKey == "" {
			printUsage()
			return
		}

		awsRegion := viper.GetString(flagConfigAwsRegion)
		if awsRegion == "" {
			printUsage()
			return
		}

		configContent, err := util.GetSecret(awsSecretKey, awsRegion)
		if err != nil {
			fmt.Printf("get aws config error, err=%s", err.Error())
			return
		}
		config = util.ParseConfigFromJson(configContent)
	} else {
		configFilePath := viper.GetString(flagConfigPath)
		if configFilePath == "" {
			printUsage()
			return
		}
		config = util.ParseConfigFromFile(configFilePath)
	}
	config.Validate()

	// init logger
	util.InitLogger(*config.LogConfig)
	util.InitAlert(config.AlertConfig)

	db, err := gorm.Open(config.DBConfig.Dialect, config.DBConfig.DBPath)
	if err != nil {
		panic(fmt.Sprintf("open db error, err=%s", err.Error()))
	}
	defer db.Close()
	model.InitTables(db)

	ascExecutor := asc.NewExecutor(config.ChainConfig.ASCProviders, config)
	ob := observer.NewObserver(db, config, ascExecutor)
	go ob.Start()

	afcExecutor, err := afc.NewExecutor(config.ChainConfig.AFCRpcAddrs, types.Network, config)
	if err != nil {
		fmt.Printf("new afc executor error, err=%s\n", err.Error())
		return
	}
	oracleRelayer := relayer.NewRelayer(db, afcExecutor, config)
	go oracleRelayer.Main()

	adm := admin.NewAdmin(config, afcExecutor)
	go adm.Serve()

	select {}
}
