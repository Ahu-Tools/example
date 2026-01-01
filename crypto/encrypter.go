package crypto

// Encrypter defines how data is secured before DB storage.
// Using an interface allows you to swap AES for RSA or a Mock for testing.
type Encrypter interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	ComputeBlindIndex(plaintext string) string
}
