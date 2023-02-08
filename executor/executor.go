package executor

import (
	"github.com/aximchain/go-sdk/common/types"
	"github.com/aximchain/go-sdk/types/msg"

	"github.com/Sotatek-huytran2/oracle-relayer/common"
)

type AfcExecutor interface {
	GetAddress() types.ValAddress
	GetCurrentSequence(chainId uint16) (int64, error)
	GetProphecy(chainId uint16, sequence int64) (*msg.Prophecy, error)

	Claim(chainId uint16, sequence uint64, payload []byte) (string, error)
}

type AscExecutor interface {
	GetBlockAndPackages(height int64) (*common.BlockAndPackageLogs, error)
}
