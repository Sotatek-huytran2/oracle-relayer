package afc

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aximchain/go-sdk/client/rpc"
	"github.com/aximchain/go-sdk/common/types"
	"github.com/aximchain/go-sdk/keys"
	"github.com/aximchain/go-sdk/types/msg"

	"github.com/Sotatek-huytran2/oracle-relayer/util"
)

type Executor struct {
	config     *util.Config
	RpcClients []rpc.Client
}

// NewExecutor returns the AFC executor instance
func NewExecutor(providers []string, network types.ChainNetwork, cfg *util.Config) (*Executor, error) {
	return &Executor{
		config:     cfg,
		RpcClients: initClients(providers, network),
	}, nil
}

// getKeyManager returns the key manager from config
func getKeyManager(config *util.ChainConfig) (keys.KeyManager, error) {
	var axcMnemonic string
	if config.AFCKeyType == util.KeyTypeAWSMnemonic {
		awsMnemonic, err := util.GetSecret(config.AFCAWSSecretName, config.AFCAWSRegion)
		if err != nil {
			return nil, err
		}
		axcMnemonic = awsMnemonic
	} else {
		axcMnemonic = config.AFCMnemonic
	}

	return keys.NewMnemonicKeyManager(axcMnemonic)
}

func initClients(providers []string, network types.ChainNetwork) []rpc.Client {
	clients := make([]rpc.Client, 0)

	for _, provider := range providers {
		client := rpc.NewRPCClient(provider, network)
		clients = append(clients, client)
	}

	return clients
}

func (e *Executor) getClient() rpc.Client {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	idx := r.Intn(len(e.RpcClients))
	return e.RpcClients[idx]
}

// GetAddress returns validator address of the oracle relayer
func (e *Executor) GetAddress() types.ValAddress {
	keyManager, err := getKeyManager(e.config.ChainConfig)
	if err != nil {
		return types.ValAddress{}
	}
	return types.ValAddress(keyManager.GetAddr())
}

// GetProphecy returns the prophecy of the given sequence
func (e *Executor) GetProphecy(chainId uint16, sequence int64) (*msg.Prophecy, error) {
	prop, err := e.getClient().GetProphecy(types.IbcChainID(chainId), sequence)
	if err != nil {
		return nil, err
	}
	return prop, err
}

// Claim sends claim to Axim Chain
func (e *Executor) Claim(chainId uint16, sequence uint64, payload []byte) (string, error) {
	client := e.getClient()

	keyManager, err := getKeyManager(e.config.ChainConfig)
	if err != nil {
		return "", fmt.Errorf("get key manager error, err=%s", err.Error())
	}
	client.SetKeyManager(keyManager)
	defer client.SetKeyManager(nil)

	res, err := client.Claim(types.IbcChainID(chainId), sequence, payload, rpc.Commit)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", fmt.Errorf("claim error, code=%d, log=%s", res.Code, res.Log)
	}
	util.Logger.Infof("claim success, tx_hash=%s", res.Hash.String())
	return res.Hash.String(), nil
}

// GetCurrentSequence return the current oracle sequence of Axim Chain
func (e *Executor) GetCurrentSequence(chainId uint16) (int64, error) {
	sequence, err := e.getClient().GetCurrentOracleSequence(types.IbcChainID(chainId))
	if err != nil {
		return 0, err
	}
	return sequence, nil
}
