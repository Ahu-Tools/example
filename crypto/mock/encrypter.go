package mock

import (
	"encoding/json"
	"errors"
)

// Envelope represents the stored data structure.
type Envelope struct {
	KeyVersion int    `json:"v"` // Which key was used?
	Nonce      []byte `json:"n"` // Random nonce for GCM
	CipherText []byte `json:"c"` // The actual encrypted data
}

type RotationManager struct {
	CurrentKeyVersion int
	Keys              map[int][]byte // Map of Version -> AES Key
	BlindIndexPepper  []byte
}

func NewRotationManager() *RotationManager {
	return &RotationManager{
		CurrentKeyVersion: 1,
		Keys: map[int][]byte{
			1: []byte("this-is-a-key"),
		},
		BlindIndexPepper: []byte("24c5833cbeebf524c839499f8b47f8886f167171347d444f71043298b1df7a12"),
	}
}

func (rm *RotationManager) Encrypt(plaintext []byte) ([]byte, error) {
	// 1. Always use the LATEST key for new writes
	key := rm.Keys[rm.CurrentKeyVersion]

	// ... (Standard AES-GCM logic here: Generate Nonce, Seal) ...
	// For brevity, we are mocking the AES logic
	cipherText := append(key, plaintext...) // MOCK encryption

	// 2. Pack into Envelope
	env := Envelope{
		KeyVersion: rm.CurrentKeyVersion,
		Nonce:      []byte("random-nonce"),
		CipherText: cipherText,
	}

	return json.Marshal(env)
}

func (rm *RotationManager) Decrypt(blob []byte) ([]byte, error) {
	var env Envelope
	if err := json.Unmarshal(blob, &env); err != nil {
		return nil, err
	}

	// 3. KEY ROTATION LOGIC:
	// We look up the key specific to THIS envelope's version.
	key, exists := rm.Keys[env.KeyVersion]
	if !exists {
		return nil, errors.New("decryption key version not found (too old?)")
	}

	// ... (Standard AES-GCM logic here: Open) ...
	// MOCK decryption
	plaintext := env.CipherText[len(key):]

	return plaintext, nil
}
