package base

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4

// Wallet stores private and public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewWallet creates and returns a new Wallet
func NewWallet() *Wallet {
	priv, _ := newKeyPair()
	return CreateWallet(&priv, &priv.PublicKey)
}

// CreateWallet initialize a wallet from the given keys
func CreateWallet(privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) *Wallet {
	// Create a wallet with the given keys, note that the PublicKey field in the
	// Wallet struct is a byte array (the concatenation of the X and Y
	// coordinates of the ecdsa.PublicKey) this is done to be easy to hash it
	// in future operations
	return &Wallet{PrivateKey: *privKey, PublicKey: pubKeyToByte(*pubKey)}
}

// GetAddress returns wallet address
// https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address
func (w Wallet) GetAddress() []byte {
	// Create a address following the logic described in the link above and
	// in the lab documentation
	pubHash := HashPubKey(w.PublicKey)
	versionedPayload := append([]byte{version}, pubHash...)
	checksum := checksum(versionedPayload)
	encoded := Base58Encode(append(versionedPayload, checksum...))

	return encoded
}

// GetStringAddress returns wallet address as string
func (w Wallet) GetStringAddress() string {
	return string(w.GetAddress())
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) []byte {
	// compute the SHA256 + RIPEMD160 hash of the pubkey
	// step 2 and 3 of: https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address
	// use the go package ripemd160:
	// https://godoc.org/golang.org/x/crypto/ripemd160
	sum := sha256.Sum256(pubKey)
	rip := ripemd160.New()
	rip.Write(sum[:])
	return rip.Sum(nil)
}

// GetPubKeyHashFromAddress returns the hash of the public key
// discarding the version and the checksum
func GetPubKeyHashFromAddress(address string) []byte {
	// Decode the address using Base58Decode and extract the hash of the pubkey
	// Look in the picture of the documentation of the lab to understand
	// how it is stored: version + pubkeyhash + checksum
	decoded := Base58Decode([]byte(address))
	return decoded[1 : len(decoded)-addressChecksumLen]
}

// ValidateAddress check if an address is valid
func ValidateAddress(address string) bool {
	// Validate a address by decoding it, extracting the
	// checksum, re-computing it using the "checksum" function
	// and comparing both.
	decoded := Base58Decode([]byte(address))
	check := decoded[len(decoded)-addressChecksumLen:]
	payload := decoded[:len(decoded)-addressChecksumLen]

	return bytes.Equal(check, checksum(payload))
}

// Checksum generates a checksum for a public key
func checksum(payload []byte) []byte {
	// Perform a double sha256 on the versioned payload
	// and return the first 4 bytes
	// Steps 5,6, and 7 of: https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address
	sum := sha256.Sum256(payload)
	sum = sha256.Sum256(sum[:])

	return sum[:addressChecksumLen]
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	// Create a new cryptographic key pair using
	// the "elliptic" and "ecdsa" package.
	// Additionally, convert the PublicKey to byte
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return *privKey, pubKeyToByte(privKey.PublicKey)
}

func pubKeyToByte(pubkey ecdsa.PublicKey) []byte {
	// step 1 of: https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses#How_to_create_Bitcoin_Address
	x := pubkey.X.Bytes()
	y := pubkey.Y.Bytes()
	return append(x, y...)
}

func encodeKeyPair(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	return encodePrivateKey(privateKey), encodePublicKey(publicKey)
}

func encodePrivateKey(privateKey *ecdsa.PrivateKey) string {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	return string(pemEncoded)
}

func encodePublicKey(publicKey *ecdsa.PublicKey) string {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncodedPub)
}

func decodeKeyPair(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	return decodePrivateKey(pemEncoded), decodePublicKey(pemEncodedPub)
}

func decodePrivateKey(pemEncoded string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	privateKey, _ := x509.ParseECPrivateKey(block.Bytes)

	return privateKey
}

func decodePublicKey(pemEncodedPub string) *ecdsa.PublicKey {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	genericPubKey, _ := x509.ParsePKIXPublicKey(blockPub.Bytes)
	publicKey := genericPubKey.(*ecdsa.PublicKey) // cast to ecdsa

	return publicKey
}
