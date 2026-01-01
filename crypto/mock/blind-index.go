package mock

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func (rm *RotationManager) ComputeBlindIndex(plaintext string) string {
	h := hmac.New(sha256.New, rm.BlindIndexPepper)
	h.Write([]byte(plaintext))
	return hex.EncodeToString(h.Sum(nil))
}
