package serialization

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"main/utils"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/pkg/errors"
)

func SerializeAndPrintForGroth16Precompile(proof groth16.Proof, vk groth16.VerifyingKey, wit witness.Witness) {
	//	----------- PUBLIC WITNESS -------------
	publicWitness, err := wit.Public()
	if err != nil {
		panic(errors.Wrap(err, "failed to extract public witness"))
	}
	buffer := &bytes.Buffer{}
	if _, err := publicWitness.WriteTo(buffer); err != nil {
		panic(errors.Wrap(err, "failed to write public witness to buffer"))
	}
	publicWitnessData := buffer.Bytes()
	fmt.Println("PUBLIC WITNESS LEN:", len(publicWitnessData))

	// 	----------- PROOF -------------
	buffer = &bytes.Buffer{}
	if _, err := proof.WriteTo(buffer); err != nil {
		panic(errors.Wrap(err, "failed to write proof to buffer"))
	}
	proofData := buffer.Bytes()
	fmt.Println("PROOF LEN:", len(proofData))

	// 	----------- VERIFIER KEY -------------
	buffer = &bytes.Buffer{}
	if _, err := vk.WriteTo(buffer); err != nil {
		panic(errors.Wrap(err, "failed to write verifier key to buffer"))
	}
	verifierKeyData := buffer.Bytes()
	fmt.Println("VERIFIER KEY LEN:", len(verifierKeyData))

	// ------------ INPUT -------------
	input := utils.Uint16ToBytes(uint16(proof.CurveID()))
	input = append(input, utils.Uint32ToBytes(uint32(len(publicWitnessData)))...)
	input = append(input, publicWitnessData...)
	input = append(input, utils.Uint32ToBytes(uint32(len(proofData)))...)
	input = append(input, proofData...)
	input = append(input, verifierKeyData...)
	fmt.Println("INPUT LEN:", len(input))
	fmt.Println("INPUT HEX: 0x" + hex.EncodeToString(input))
}
