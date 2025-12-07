//go:build js

package js

import (
	"fmt"
	"syscall/js"
	"wasm-test/model"
)

const sdkName = "secubird"

var jsSdk js.Value

var currentUser *model.User = nil

func init() {
	jsSdk = js.ValueOf(make(map[string]any))
	export()
	js.Global().Set(sdkName, jsSdk)
}

func export() {
	jsSdk.Set("loginOrRegister", js.FuncOf(loginOrRegister))
	jsSdk.Set("logout", js.FuncOf(logout))
	jsSdk.Set("encryptMessage", js.FuncOf(encryptMessage))
	jsSdk.Set("decryptMessage", js.FuncOf(decryptMessage))
}

func loginOrRegister(this js.Value, args []js.Value) any {
	if currentUser != nil {
		js.Global().Call("eval", "console.log('Already logged in!')")
		return js.Undefined()
	}

	if len(args) < 2 {
		js.Global().Call("eval", "console.log('Email and password required!')")
		return js.Undefined()
	}

	if args[0].Type() != js.TypeString || args[1].Type() != js.TypeString {
		js.Global().Call("eval", "console.log('Email and password need to be strings!')")
		return js.Undefined()
	}

	email, password := args[0].String(), args[1].String()
	user, err := model.NewUser(email, password)
	if err != nil {
		js.Global().Call("eval", fmt.Sprintf("console.log('%s')", err))
	}

	currentUser = &user
	js.Global().Call("eval", "console.log('Successfully logged in or registered!')")
	fmt.Println(currentUser)
	return js.Undefined()
}

func logout(this js.Value, args []js.Value) any {
	if currentUser == nil {
		js.Global().Call("eval", "Not logged in!")
		return js.Undefined()
	}

	currentUser = nil
	js.Global().Call("eval", "console.log('Successfully logged out!')")
	return js.Undefined()
}

func encryptMessage(this js.Value, args []js.Value) any {
	if currentUser == nil {
		js.Global().Call("eval", "console.log('You need to be logged in!')")
		return js.Undefined()
	}

	if len(args) < 1 {
		js.Global().Call("eval", "console.log('Message required!')")
		return js.Undefined()
	}

	if args[0].Type() != js.TypeString {
		js.Global().Call("eval", "console.log('Message needs to be a string!')")
		return js.Undefined()
	}

	message := args[0].String()
	cipherText := currentUser.EncryptMessage(message)
	return js.ValueOf(cipherText)
}

func decryptMessage(this js.Value, args []js.Value) any {
	if currentUser == nil {
		js.Global().Call("eval", "console.log('You need to be logged in!')")
		return js.Undefined()
	}

	if len(args) < 1 {
		js.Global().Call("eval", "console.log('Message required!')")
		return js.Undefined()
	}

	if args[0].Type() != js.TypeString {
		js.Global().Call("eval", "console.log('Message needs to be a string!')")
		return js.Undefined()
	}

	message := args[0].String()
	plainText, err := currentUser.DecryptMessage(message)
	if err != nil {
		js.Global().Call("eval", fmt.Sprintf("console.log('%s')", err))
		return js.Undefined()
	}

	return js.ValueOf(plainText)
}
