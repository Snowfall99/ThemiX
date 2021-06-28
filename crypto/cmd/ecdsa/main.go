// Copyright 2021 The themix authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"go.themix.io/crypto/ecdsa"
	"go.themix.io/crypto/sha256"
)

func main() {
	_, err := ecdsa.GenerateEcdsaKey("./")

	if err != nil {
		fmt.Println("Failed to generate a key: ", err)
		return
	}

	priv, _ := ecdsa.LoadKey("./")

	hash, err := sha256.ComputeHash([]byte{'a', 'b', 'c'})

	if err != nil {
		fmt.Println("Failed to compute hash of a message: ", err)
	}

	sig, err := ecdsa.SignECDSA(priv, hash)

	if err != nil {
		fmt.Println("Failed to sign a mesage: ", err)
	}

	b, err := ecdsa.VerifyECDSA(&priv.PublicKey, sig, hash)

	if err != nil {
		fmt.Println("Failed to verify a message: ", err)
	}

	if !b {
		fmt.Println("Message verification failed")
	} else {
		fmt.Println("Message verification succeeded")
	}
}
