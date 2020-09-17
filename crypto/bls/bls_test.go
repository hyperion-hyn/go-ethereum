package bls

import (
	"strings"
	"testing"
	"time"

	"github.com/hyperion-hyn/bls/ffi/go/bls"

	"github.com/ethereum/go-ethereum/crypto"
)

// Test the basic functionality of a BLS multi-sig mask.
func TestNewMask(test *testing.T) {
	pubKey1 := RandPrivateKey().GetPublicKey()
	pubKey2 := RandPrivateKey().GetPublicKey()
	pubKey3 := RandPrivateKey().GetPublicKey()

	mask, err := NewMask([]*bls.PublicKey{pubKey1, pubKey2, pubKey3}, pubKey1)

	if err != nil {
		test.Errorf("Failed to create a new Mask: %s", err)
	}

	if mask.Len() != 1 {
		test.Errorf("Mask created with wrong size: %d", mask.Len())
	}

	enabled, err := mask.KeyEnabled(pubKey1)
	if !enabled || err != nil {
		test.Errorf("My key pubKey1 should have been enabled: %s", err)
	}

	if mask.CountEnabled() != 1 {
		test.Error("Only one key should have been enabled")
	}

	if mask.CountTotal() != 3 {
		test.Error("Should have a total of 3 keys")
	}
}

func TestNewMaskWithAbsentPublicKey(test *testing.T) {
	pubKey1 := RandPrivateKey().GetPublicKey()
	pubKey2 := RandPrivateKey().GetPublicKey()
	pubKey3 := RandPrivateKey().GetPublicKey()
	pubKey4 := RandPrivateKey().GetPublicKey()

	mask, err := NewMask([]*bls.PublicKey{pubKey1, pubKey2, pubKey3}, pubKey4)

	if err == nil {
		test.Errorf("Failed to create a new Mask: %s", err)
	}

	if mask != nil {
		test.Errorf("Expected failure to create a new mask")
	}

}

func TestThreshHoldPolicy(test *testing.T) {
	pubKey1 := RandPrivateKey().GetPublicKey()
	pubKey2 := RandPrivateKey().GetPublicKey()
	pubKey3 := RandPrivateKey().GetPublicKey()

	mask, err := NewMask([]*bls.PublicKey{pubKey1, pubKey2, pubKey3}, pubKey1)

	if err != nil {
		test.Errorf("Failed to create a new Mask: %s", err)
	}

	if mask.Len() != 1 {
		test.Errorf("Mask created with wrong size: %d", mask.Len())
	}

	threshHoldPolicy := *NewThresholdPolicy(1)

	mask.SetKey(pubKey1, true)
	mask.SetKey(pubKey2, true)

	if mask.CountEnabled() != 2 {
		test.Errorf("Number of enabled nodes: %d , expected count = 2 ", mask.CountEnabled())
	}

	if !threshHoldPolicy.Check(mask) {
		test.Error("Number of enabled nodes less than threshold")
	}

	mask.SetKey(pubKey1, false)
	mask.SetKey(pubKey2, false)

	if threshHoldPolicy.Check(mask) {
		test.Error("Number of enabled nodes more than equal to threshold")
	}
}

func TestCompletePolicy(test *testing.T) {
	pubKey1 := RandPrivateKey().GetPublicKey()
	pubKey2 := RandPrivateKey().GetPublicKey()
	pubKey3 := RandPrivateKey().GetPublicKey()

	mask, err := NewMask([]*bls.PublicKey{pubKey1, pubKey2, pubKey3}, pubKey1)

	if err != nil {
		test.Errorf("Failed to create a new Mask: %s", err)
	}

	if mask.Len() != 1 {
		test.Errorf("Mask created with wrong size: %d", mask.Len())
	}

	completePolicy := CompletePolicy{}

	mask.SetKey(pubKey1, true)
	mask.SetKey(pubKey2, true)
	mask.SetKey(pubKey3, true)

	if mask.CountEnabled() != 3 {
		test.Errorf("Number of enabled nodes: %d , expected count = 3 ", mask.CountEnabled())
	}

	if !completePolicy.Check(mask) {
		test.Error("Number of enabled nodes not equal to total count")
	}

	mask.SetKey(pubKey1, false)

	if completePolicy.Check(mask) {
		test.Error("Number of enabled nodes equal to total count")
	}
}

func TestAggregatedSignature(test *testing.T) {
	var sec bls.SecretKey
	sec.SetByCSPRNG()

	signs := []*bls.Sign{sec.Sign("message1"), sec.Sign("message2")}

	multiSignature := AggregateSig(signs)

	str := multiSignature.SerializeToHexStr()

	if strings.Compare(multiSignature.SerializeToHexStr(), "0") == 0 {
		test.Error("Error creating multisignature", str)
	}
}

func TestAggregateMasks(test *testing.T) {
	message := []byte("message")
	newMessage := []byte("message")
	emptyMessage := []byte("")

	aggMask, err := AggregateMasks(message, emptyMessage)
	if aggMask != nil {
		test.Error("Expected mismatching bitmap lengths")
	}
	if err == nil {
		test.Error("Expected error thrown because of bitmap length mismatch")
	}

	if _, err := AggregateMasks(message, newMessage); err != nil {
		test.Error("Error thrown in aggregating masks")
	}
}

func TestEnableKeyFunctions(test *testing.T) {
	pubKey1 := RandPrivateKey().GetPublicKey()
	pubKey2 := RandPrivateKey().GetPublicKey()
	pubKey3 := RandPrivateKey().GetPublicKey()
	pubKey4 := RandPrivateKey().GetPublicKey()

	mask, err := NewMask([]*bls.PublicKey{pubKey1, pubKey2, pubKey3}, pubKey1)

	if err != nil {
		test.Errorf("Failed to create a new Mask: %s", err)
	}
	length := mask.Len()
	_ = length

	if mask.Len() != 1 {
		test.Errorf("Mask created with wrong size: %d", mask.Len())
	}

	mask.SetBit(0, true)
	mask.SetBit(1, false)
	mask.SetBit(0, true)

	if err := mask.SetBit(5, true); err == nil {
		test.Error("Expected index out of range error")
	}

	enabledKeys := mask.GetPubKeyFromMask(true)
	disabledKeys := mask.GetPubKeyFromMask(false)

	if len(enabledKeys) != 1 {
		test.Error("Count of enabled keys doesn't match")
	}

	if len(disabledKeys) != 2 {
		test.Error("Count of disabled keys don't match")
	}

	if _, error := mask.KeyEnabled(pubKey4); error == nil {
		test.Error("Expected key not found error")
	}

	if _, error := mask.IndexEnabled(5); error == nil {
		test.Error("Expected index out of range error")
	}

	if err := mask.SetKey(pubKey4, true); err == nil {
		test.Error("Expected key nout found error")
	}
}

func TestCopyParticipatingMask(test *testing.T) {
	pubKey1 := RandPrivateKey().GetPublicKey()
	pubKey2 := RandPrivateKey().GetPublicKey()

	mask, _ := NewMask([]*bls.PublicKey{pubKey1, pubKey2}, pubKey1)

	clonedMask := mask.Mask()

	if len(clonedMask) != 1 {
		test.Error("Length of clonedMask doesn't match")
	}

}

func TestSetMask(test *testing.T) {
	pubKey1 := RandPrivateKey().GetPublicKey()
	pubKey2 := RandPrivateKey().GetPublicKey()

	mask, _ := NewMask([]*bls.PublicKey{pubKey1, pubKey2}, pubKey1)

	_ = mask
	maskBytes := []byte{3}
	mask.SetMask(maskBytes)

	if mask.CountEnabled() != 2 {
		test.Error("Count of Enabled nodes doesn't match")
	}

	newMaskBytes := []byte{3, 2}

	if err := mask.SetMask(newMaskBytes); err == nil {
		test.Error("Expected mismatching Bitmap lengths")
	}
}

func TestSignVerify(t *testing.T) {
	var signature []byte
	secretKey, err := crypto.GenerateBLSKey()
	if err != nil {
		t.Errorf("failed to generate bls key: %v", err)
	}

	data := time.Now().String()
	hash := crypto.Keccak256Hash([]byte(data)).Bytes()

	{
		sign := secretKey.SignHash(hash)
		signature = sign.Serialize()
	}

	{
		var sign bls.Sign
		err := sign.Deserialize(signature)
		if err != nil {
			t.Errorf("failed to deserialize signature: %v", err)
		}

		publicKey := secretKey.GetPublicKey()
		if ok := sign.VerifyHash(publicKey, hash); !ok {
			t.Errorf("failed to verify hash")
		}
	}
}

func TestMultipleSign(t *testing.T) {
	const KEY_COUNT = 10
	var keys [KEY_COUNT]*bls.SecretKey
	for i := 0; i < KEY_COUNT; i++ {
		keys[i] = RandPrivateKey()
	}

	data := time.Now().String()
	hash := crypto.Keccak256Hash([]byte(data)).Bytes()

	var aggregatedPublicKey bls.PublicKey
	var aggregatedSign bls.Sign
	for i := 0; i < KEY_COUNT-1; i++ {
		aggregatedPublicKey.Add(keys[i].GetPublicKey())
		sign := keys[i].SignHash(hash)
		aggregatedSign.Add(sign)
	}

	if ok := aggregatedSign.VerifyHash(&aggregatedPublicKey, hash); !ok {
		t.Errorf("failed to verify aggregated signatures.")
	}

	{
		if ok := aggregatedSign.VerifyHash(keys[0].GetPublicKey(), hash); ok {
			t.Errorf("failed to verify aggregated signatures.")
		}
	}

	{
		sign := keys[1].SignHash(hash)
		if ok := sign.VerifyHash(keys[0].GetPublicKey(), hash); ok {
			t.Errorf("failed to verify aggregated signatures.")
		}
	}

	{
		// repeat
		aggregatedPublicKey.Add(keys[0].GetPublicKey())
		if ok := aggregatedSign.VerifyHash(&aggregatedPublicKey, hash); ok {
			t.Errorf("failed to verify aggregated signatures.")
		}
		aggregatedPublicKey.Sub(keys[0].GetPublicKey())
	}

	{
		// additional
		aggregatedPublicKey.Add(keys[KEY_COUNT-1].GetPublicKey())
		if ok := aggregatedSign.VerifyHash(&aggregatedPublicKey, hash); ok {
			t.Errorf("failed to verify aggregated signatures.")
		}
		aggregatedPublicKey.Sub(keys[KEY_COUNT-1].GetPublicKey())
	}

	{
		// removal
		aggregatedPublicKey.Sub(keys[0].GetPublicKey())
		if ok := aggregatedSign.VerifyHash(&aggregatedPublicKey, hash); ok {
			t.Errorf("failed to verify aggregated signatures.")
		}
		aggregatedPublicKey.Add(keys[0].GetPublicKey())
	}

	if ok := aggregatedSign.VerifyHash(&aggregatedPublicKey, hash); !ok {
		t.Errorf("failed to verify aggregated signatures.")
	}
}
