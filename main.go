package main

import (
	"bytes"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/pkg/errors"
	"main/circuits"
)

func main() {
	var (
		ecID = ecc.BN254
	)

	circuit, assignment := circuits.SampleZKMessage()

	// ---------------- COMPILE CIRCUIT -------------------
	ccs, err := frontend.Compile(ecID.ScalarField(), r1cs.NewBuilder, circuit)
	if err != nil {
		panic(errors.Wrap(err, "failed to compile circuit"))
	}

	// ---------------- SETUP PK AND VK --------------------
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		panic(errors.Wrap(err, "failed to setup pk and vk"))
	}

	// ---------------- DEFINE WITNESS ----------
	wit, err := frontend.NewWitness(assignment, ecID.ScalarField())
	if err != nil {
		panic(errors.Wrap(err, "failed to instantiate new witness"))
	}

	// ---------------- MAKE PROOF -------------------
	proof, err := groth16.Prove(ccs, pk, wit)
	if err != nil {
		panic(errors.Wrap(err, "failed to prove"))
	}

	// ----------------- VERIFY PROOF --------------------
	publicWitness, err := wit.Public()
	if err != nil {
		panic(errors.Wrap(err, "failed to extract public witness"))
	}
	if err := groth16.Verify(proof, vk, publicWitness); err != nil {
		panic(errors.Wrap(err, "failed to verify proof"))
	}

	buffer := new(bytes.Buffer)

	if _, err := proof.WriteTo(buffer); err != nil {
		panic(errors.Wrap(err, "failed to write proof to buffer"))
	}
	proofData := buffer.Bytes()

	buffer.Reset()

	if _, err := vk.WriteTo(buffer); err != nil {
		panic(errors.Wrap(err, "failed to write vk to buffer"))
	}
	verifierKeyData := buffer.Bytes()

	buffer.Reset()

	if _, err := publicWitness.WriteTo(buffer); err != nil {
		panic(errors.Wrap(err, "failed to write public witness to buffer"))
	}
	publicWitnessData := buffer.Bytes()

	if err := groth16.Verify(proof, vk, publicWitness); err != nil {
		panic(errors.Wrap(err, "failed to verify proof | serialized"))
	}

	fmt.Println(proofData, verifierKeyData, publicWitnessData)
}
