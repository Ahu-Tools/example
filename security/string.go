package security

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/Ahu-Tools/example/crypto"
)

// SecureString wraps a string that should be encrypted at rest.
type SecureString string

func SecureStringPtr(s string) *SecureString {
	if len(s) == 0 {
		return nil
	}
	secure := SecureString(s)
	return &secure
}

// Value is called AUTOMATICALLY by the database driver when writing to DB.
// It converts the simple string -> Encrypted Byte Slice.
func (s SecureString) Value() (driver.Value, error) {
	if s == "" {
		return nil, nil
	}

	// Call your global or singleton crypto service
	encryptedBytes, err := crypto.GlobalEncrypter.Encrypt([]byte(s))
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}

	return encryptedBytes, nil
}

// Scan is called AUTOMATICALLY by the database driver when reading from DB.
// It converts the Encrypted Byte Slice -> simple string.
func (s *SecureString) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	decryptedBytes, err := crypto.GlobalEncrypter.Decrypt(bytes)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	*s = SecureString(decryptedBytes)
	return nil
}

// MarshalJSON ensures we never accidentally log or return the raw secret in JSON.
func (s SecureString) MarshalJSON() ([]byte, error) {
	// Return a masked value or null
	return []byte(`"***REDACTED***"`), nil
}
