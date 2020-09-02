package microstaking

import (
	common2 "github.com/ethereum/go-ethereum/staking/types/common"
	"github.com/pkg/errors"
	"testing"
)

var (
	errBLSKeysNotMatchSigs = errors.New("bls keys and corresponding signatures could not be verified")
)

func TestVerifyBLSKey(t *testing.T) {
	pairs := makeBLSPubSigPairs(5)
	tests := []struct {
		pubIndex int
		sigIndex int
		expErr     error
	}{
		{0, 0, nil},
		{0, 1, errBLSKeysNotMatchSigs},
	}
	for i, test := range tests {
		pubs := getPubsFromPairs(pairs, []int{test.pubIndex})
		sigs := getSigsFromPairs(pairs, []int{test.sigIndex})

		err := common2.VerifyBLSKey(pubs.Keys[0], &sigs[0])
		if assErr := assertError(err, test.expErr); assErr != nil {
			t.Errorf("Test %v: %v", i, assErr)
		}
	}
}
