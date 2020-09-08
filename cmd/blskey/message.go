// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/hyperion-hyn/bls/ffi/go/bls"
	"gopkg.in/urfave/cli.v1"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
)

type outputSign struct {
	Signature string
}

var msgfileFlag = cli.StringFlag{
	Name:  "msgfile",
	Usage: "file containing the message to sign/verify",
}

var hashFlag = cli.StringFlag{
	Name:  "hash",
	Usage: "hash of message to sign/verify",
}

var commandSignMessage = cli.Command{
	Name:      "signmessage",
	Usage:     "sign a message",
	ArgsUsage: "<keyfile> <message>",
	Description: `
Sign the message with a keyfile.

To sign a message contained in a file, use the --msgfile flag.
To sign a hash of message, use the --hash flag.
`,
	Flags: []cli.Flag{
		passphraseFlag,
		jsonFlag,
		msgfileFlag,
		hashFlag,
	},
	Action: func(ctx *cli.Context) error {
		message, isMessage := getMessage(ctx, 1)

		// Load the keyfile.
		keyfilepath := ctx.Args().First()
		keyjson, err := ioutil.ReadFile(keyfilepath)
		if err != nil {
			utils.Fatalf("Failed to read the keyfile at '%s': %v", keyfilepath, err)
		}

		// Decrypt key with passphrase.
		passphrase := getPassphrase(ctx)
		key, err := keystore.DecryptBLSKey(keyjson, passphrase)
		if err != nil {
			utils.Fatalf("Error decrypting key: %v", err)
		}

		var signature *bls.Sign
		if isMessage {
			signature = key.PrivateKey.SignHash(signHash(message))
		} else {
			signature = key.PrivateKey.SignHash(common.HexToHash(string(message)).Bytes())
		}
		if signature == nil {
			utils.Fatalf("Failed to sign message: %v", err)
		}
		out := outputSign{Signature: hex.EncodeToString(signature.Serialize())}
		if ctx.Bool(jsonFlag.Name) {
			mustPrintJSON(out)
		} else {
			fmt.Println("Signature:", out.Signature)
		}
		return nil
	},
}

type outputVerify struct {
	Success   bool
	PublicKey string
}

var commandVerifyMessage = cli.Command{
	Name:      "verifymessage",
	Usage:     "verify the signature of a signed message",
	ArgsUsage: "<public-key> <signature> <message>",
	Description: `
Verify the signature of the message.
It is possible to refer to a file containing the message.`,
	Flags: []cli.Flag{
		jsonFlag,
		msgfileFlag,
		hashFlag,
	},
	Action: func(ctx *cli.Context) error {
		pubKeyHex := ctx.Args().First()
		signatureHex := ctx.Args().Get(1)
		message, isMessage := getMessage(ctx, 2)

		var publicKey bls.PublicKey
		data, err := hex.DecodeString(pubKeyHex)
		if err != nil {
			utils.Fatalf("Public Key encoding is not hexadecimal: %v", err)
		}
		err = publicKey.Deserialize(data)
		if err != nil {
			utils.Fatalf("Public Key is not deserialized: %v", err)
		}

		var signature bls.Sign
		data, err = hex.DecodeString(signatureHex)
		if err != nil {
			utils.Fatalf("Signature encoding is not hexadecimal: %v", err)
		}
		err = signature.Deserialize(data)
		if err != nil {
			utils.Fatalf("Signature is not deserialized: %v", err)
		}

		var success bool
		if isMessage {
			success = signature.VerifyHash(&publicKey, signHash(message))
		} else {
			success = signature.VerifyHash(&publicKey, common.HexToHash(string(message)).Bytes())
		}
		if !success {
			utils.Fatalf("Signature verification failed")
		}

		out := outputVerify{
			Success:   success,
			PublicKey: hex.EncodeToString(publicKey.Serialize()),
		}
		if ctx.Bool(jsonFlag.Name) {
			mustPrintJSON(out)
		} else {
			if out.Success {
				fmt.Println("Signature verification successful!")
			} else {
				fmt.Println("Signature verification failed!")
			}
			fmt.Println("PublicKey:", out.PublicKey)
		}
		return nil
	},
}

func getMessage(ctx *cli.Context, msgarg int) (data []byte, isMessage bool) {
	if hash := ctx.String("hash"); hash != "" {
		if len(ctx.Args()) > msgarg {
			utils.Fatalf("Can't use --hash and message argument at the same time.")
		}
		return []byte(hash), false
	} else if file := ctx.String("msgfile"); file != "" {
		if len(ctx.Args()) > msgarg {
			utils.Fatalf("Can't use --msgfile and message argument at the same time.")
		}
		msg, err := ioutil.ReadFile(file)
		if err != nil {
			utils.Fatalf("Can't read message file: %v", err)
		}
		return msg, true
	} else if len(ctx.Args()) == msgarg+1 {
		return []byte(ctx.Args().Get(msgarg)), true
	}
	utils.Fatalf("Invalid number of arguments: want %d, got %d", msgarg+1, len(ctx.Args()))
	return nil, false
}
