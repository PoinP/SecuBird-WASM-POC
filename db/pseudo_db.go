package db

import "slices"

type UserRow struct {
	Email     string
	PublicKey []byte
	Modulus   []byte

	Iv   []byte
	Salt []byte

	EPrimes     [][]byte
	EPrivateKey []byte
}

var Users []UserRow
var Messages []MessageRow

func AddUsers(u UserRow) {
	Users = append(Users, u)
}

func GetUser(email string) (UserRow, bool) {
	userIdx := slices.IndexFunc(Users, func(u UserRow) bool { return u.Email == email })
	if userIdx == -1 {
		return UserRow{}, false
	}

	return Users[userIdx], true
}

type MessageRow struct {
	Iv            []byte
	CipherMessage string
}

func AddMessage(cipherMessage string, iv []byte) {
	Messages = append(Messages, MessageRow{iv, cipherMessage})
}

func GetMessage(cipherMessage string) (MessageRow, bool) {
	messageIdx := slices.IndexFunc(Messages, func(m MessageRow) bool { return m.CipherMessage == cipherMessage })
	if messageIdx == -1 {
		return MessageRow{}, false
	}

	return Messages[messageIdx], true
}
