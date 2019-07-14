package types
import (
	//"bytes"
	//"errors"
	//"fmt"
	"time"

	//"github.com/tendermint/tendermint/crypto"
	//cmn "github.com/tendermint/tendermint/libs/common"
)
type ChangeMessage struct {
	Type             SignedMsgType `json:"type"`
	Height           int64         `json:"height"`
	Round            int           `json:"round"`
	Timestamp        time.Time     `json:"timestamp"`
	ValidatorAddress Address       `json:"validator_address"`
	ValidatorIndex   int           `json:"validator_index"`
	Signature        []byte        `json:"signature"`
}
