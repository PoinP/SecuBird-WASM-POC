package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"wasm-test/db"

	"golang.org/x/crypto/argon2"
)

// Note: Resuing IVs is bad. Since this is a POC for a WASM build, it is easier this way
// Best approach perhaps would be saving the IV at the front of the ciphertexts

type CipherKey []byte

type User struct {
	email     string
	accessKey *rsa.PrivateKey
}

func NewUser(email, password string) (User, error) {
	user, ok := db.GetUser(email)
	if !ok {
		return registerUser(email, password)
	}

	return loginUser(user, password)
}

func (u User) EncryptMessage(message string) string {
	h := sha256.New()
	h.Write(binary.BigEndian.AppendUint64(make([]byte, 0, 8), uint64(u.accessKey.E)))
	key := h.Sum(nil)

	iv := make([]byte, 16)
	rand.Read(iv)

	block, _ := aes.NewCipher(key)
	blockCipher := cipher.NewCTR(block, iv)

	cipherMessage := make([]byte, len(message))
	blockCipher.XORKeyStream(cipherMessage, []byte(message))

	cipherText := hex.EncodeToString(cipherMessage)
	db.AddMessage(cipherText, iv)

	return cipherText
}

func (u User) DecryptMessage(message string) (string, error) {
	cipherMessage, ok := db.GetMessage(message)
	if !ok {
		return "", fmt.Errorf("user: message not found in system")
	}

	cipherText, err := hex.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("user: malformed message: %v", err)
	}

	h := sha256.New()
	h.Write(binary.BigEndian.AppendUint64(make([]byte, 0, 8), uint64(u.accessKey.E)))
	key := h.Sum(nil)

	block, _ := aes.NewCipher(key)
	blockCipher := cipher.NewCTR(block, cipherMessage.Iv)

	plainText := make([]byte, len(cipherText))
	blockCipher.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}

func registerUser(email, password string) (User, error) {
	salt := make([]byte, 32)
	rand.Read(salt)

	masterKey := argon2.IDKey([]byte(password), salt, 2, 64*1024, 4, 32)

	rsaKeys, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return User{}, fmt.Errorf("user: an error occured when registering user")
	}

	iv := make([]byte, 16)
	rand.Read(iv)

	block, _ := aes.NewCipher(masterKey)
	blockCipher := cipher.NewCTR(block, iv)

	privateKey := rsaKeys.D.Bytes()
	ePrivateKey := make([]byte, len(privateKey))
	blockCipher.XORKeyStream(ePrivateKey, privateKey)

	primes := [2][]byte{rsaKeys.Primes[0].Bytes(), rsaKeys.Primes[1].Bytes()}
	ePrimes := xorStreamPrimes(primes, masterKey, iv)

	userRow := db.UserRow{
		Email:       email,
		PublicKey:   binary.BigEndian.AppendUint64(make([]byte, 0, 8), uint64(rsaKeys.E)),
		Modulus:     rsaKeys.N.Bytes(),
		Iv:          iv,
		Salt:        salt,
		EPrimes:     ePrimes,
		EPrivateKey: ePrivateKey,
	}

	db.AddUsers(userRow)

	return User{
		email:     email,
		accessKey: rsaKeys,
	}, nil
}

func loginUser(user db.UserRow, password string) (User, error) {
	masterKey := argon2.IDKey([]byte(password), user.Salt, 2, 64*1024, 4, 32)

	block, _ := aes.NewCipher(masterKey)
	blockCipher := cipher.NewCTR(block, user.Iv)

	privateKey := make([]byte, len(user.EPrivateKey))
	blockCipher.XORKeyStream(privateKey, user.EPrivateKey)

	primes := xorStreamPrimes([2][]byte(user.EPrimes), masterKey, user.Iv)

	rsaKeys := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: new(big.Int).SetBytes(user.Modulus),
			E: int(binary.BigEndian.Uint64(user.PublicKey)),
		},
		D: new(big.Int).SetBytes(privateKey),
		Primes: []*big.Int{
			new(big.Int).SetBytes(primes[0]),
			new(big.Int).SetBytes(primes[1]),
		},
	}

	rsaKeys.Precompute()

	err := rsaKeys.Validate()
	if err != nil {
		return User{}, err
	}

	return User{
		email:     user.Email,
		accessKey: rsaKeys,
	}, nil
}

func xorStreamPrimes(primes [2][]byte, masterKey, iv []byte) [][]byte {
	block, _ := aes.NewCipher(masterKey)

	ePrimes := [2][]byte{}
	for i, p := range primes {
		blockCipher := cipher.NewCTR(block, iv)
		ePrimes[i] = make([]byte, len(p))
		blockCipher.XORKeyStream(ePrimes[i], p)
	}

	return ePrimes[:]
}
