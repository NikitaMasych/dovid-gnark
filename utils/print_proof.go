package utils

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/pkg/errors"
)

func PrintGroth16BN254Proof(proof groth16.Proof) {
	const fpSize = 32
	buffer := &bytes.Buffer{}
	if _, err := proof.WriteRawTo(buffer); err != nil {
		panic(errors.Wrap(err, "failed to write proof to buffer"))
	}
	proofData := buffer.Bytes()

	var a [2]*big.Int
	a[0] = big.NewInt(0).SetBytes(proofData[fpSize*0 : fpSize*1])
	a[1] = big.NewInt(0).SetBytes(proofData[fpSize*1 : fpSize*2])

	var b [2][2]*big.Int
	b[0][0] = big.NewInt(0).SetBytes(proofData[fpSize*2 : fpSize*3])
	b[0][1] = big.NewInt(0).SetBytes(proofData[fpSize*3 : fpSize*4])
	b[1][0] = big.NewInt(0).SetBytes(proofData[fpSize*4 : fpSize*5])
	b[1][1] = big.NewInt(0).SetBytes(proofData[fpSize*5 : fpSize*6])

	var c [2]*big.Int
	c[0] = big.NewInt(0).SetBytes(proofData[fpSize*6 : fpSize*7])
	c[1] = big.NewInt(0).SetBytes(proofData[fpSize*7 : fpSize*8])

	fmt.Printf("%v\n", a)
	fmt.Printf("%v\n", b)
	fmt.Printf("%v\n", c)
}
