package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/pkg/errors"
	"os"
)

// CubicCircuit defines a simple circuit
// x**3 + x + 5 == y
type CubicCircuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
// x**3 + x + 5 == y
func (circuit *CubicCircuit) Define(api frontend.API) error {
	x3 := api.Mul(circuit.X, circuit.X, circuit.X)
	api.AssertIsEqual(circuit.Y, api.Add(x3, circuit.X, 5))
	return nil
}

func main() {
	// compiles our circuit into a R1CS
	var circuit CubicCircuit
	ccs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)

	// groth16 zkSNARK: Setup
	pk, vk, _ := groth16.Setup(ccs)

	const verifySolidityPath = "artifacts/verifier_groth16.sol"
	f, _ := os.OpenFile(verifySolidityPath, os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	vk.ExportSolidity(f)

	// witness definition
	assignment := CubicCircuit{X: 3, Y: 35}
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	// groth16: Prove & Verify

	proof, _ := groth16.Prove(ccs, pk, witness)
	if err := groth16.Verify(proof, vk, publicWitness); err != nil {
		panic(errors.Wrap(err, "failed to verify proof"))
	}

	// 	----------- PROOF -------------
	buffer := &bytes.Buffer{}
	if _, err := proof.WriteTo(buffer); err != nil {
		panic(errors.Wrap(err, "failed to write proof to buffer"))
	}
	proofData := buffer.Bytes()
	fmt.Println("PROOF LEN: ", len(proofData))
	fmt.Println(proofData)

	//	----------- PUBLIC WITNESS -------------
	buffer = &bytes.Buffer{}
	if _, err := publicWitness.WriteTo(buffer); err != nil {
		panic(errors.Wrap(err, "failed to write public witness to buffer"))
	}
	publicWitnessData := buffer.Bytes()
	fmt.Println("PUBLIC WITNESS LEN: ", len(publicWitnessData))
	fmt.Println(publicWitnessData)

	// 	----------- VERIFIER KEY -------------
	buffer = &bytes.Buffer{}
	if _, err := vk.WriteTo(buffer); err != nil {
		panic(errors.Wrap(err, "failed to write verifier key to buffer"))
	}
	verifierKeyData := buffer.Bytes()
	fmt.Println("VERIFIER KEY LEN: ", len(verifierKeyData))
	fmt.Println(verifierKeyData)

	// "length of the proof" + proof + "length of the witness" + witness + verifier_key
	input := intToBytes(uint32(len(proofData)))
	input = append(input, proofData...)
	input = append(input, intToBytes(uint32(len(publicWitnessData)))...)
	input = append(input, publicWitnessData...)
	input = append(input, verifierKeyData...)

	fmt.Println("\n\n\n\nINPUT")
	fmt.Println(input)

	hexInput := hex.EncodeToString(input)
	fmt.Println(hexInput)
}

func intToBytes(number uint32) []byte {
	byteArray := make([]byte, 4)
	binary.BigEndian.PutUint32(byteArray, number)
	return byteArray
}

// [0 0 0 1 0 0 0 1 0 0 0 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 35 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 3]
// [0 0 0 1 0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 35]

// UINT32 Number of public + UINT32 Number of private + UINT32 Number of field elements in array + VECTOR of FIELD ELEMENTS
