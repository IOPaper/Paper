package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/IOPaper/Paper/crypto"
	"github.com/IOPaper/Paper/utils"
)

var (
	output string
)

func init() {
	flag.StringVar(&output, "output", "./", "secp256k1 public key and private key output dir")
}

func main() {
	flag.Parse()
	if output == "" {
		fmt.Println("must specify an output directory")
		return
	}
	secp256k1, err := crypto.NewSecp256k1()
	if err != nil {
		fmt.Println("generate secp256k1 keypair failed error:", err)
		return
	}
	fmt.Println(">>>>> start generate keypair <<<<<")
	pri, pub := secp256k1.ExportKeypair()
	fmt.Println("===================================================")
	fmt.Printf(
		"private key fingerprint: %s\npublic key fingerprint: %s\n",
		crypto.NewFingerprint(bytes.NewBuffer(pri)).Hex(),
		crypto.NewFingerprint(bytes.NewBuffer(pub)).Hex(),
	)
	fmt.Println("===================================================")
	if err = utils.Write(output+"/pubkey", bytes.NewBuffer(pub)); err != nil {
		fmt.Println("write public key failed, error:", err)
		return
	}
	if err = utils.Write(output+"/prikey", bytes.NewBuffer(pri)); err != nil {
		fmt.Println("write private key failed, error:", err)
		return
	}
	fmt.Println(">>>>> DONE <<<<<")
}
