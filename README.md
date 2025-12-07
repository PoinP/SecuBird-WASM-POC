# SecuBird JS WASM POC

## How to run the app:
> Prerequisites: golang 1.25.5 and make

1. Open your terminal
2. Use `make build` to build the WASM binary and copy the `.js` dependencies
3. A new folder named `bin` should appear with the binaries
4. Use `make serve` to start a local server where the WASM binary will run.
5. Server should be running at `http://localhost:8080`
6. When you open the web page a single button with a label `Run` should appear
7. After clicking `Run` in the console of your browser should appear `SDK is ready to be used!` (Note: Subsequent clicks will break the SDK)
8. The WASM module is running and you can play around with it

## Module API
In order to access the module, an object inside the JS Environment name `secubird` is created.
This object will be used for communication with the Go code inside the browser.

- `secubird.loginOrRegister(string, string)` -> registers and/or logs in a user into their account
- `secubird.logout()` -> logs out a user
- `secubird.encryptMessage(string)` -> encrypts a message. returns the encrypted message
- `secubird.decryptMessage(string)` -> decrypts a message. returns the plaintext of the encrypted message

> Note: Encryption and Decryption works only if a user has been logged into their account

> Note: Decryptions happens only on messages encrypted through a call to the `encryptMessage` method

> Note: The Go code has a basic in-memory database for test purposes. The stored data gets cleared on page refresh