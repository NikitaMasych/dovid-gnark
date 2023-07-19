package utils

import (
	"os"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/pkg/errors"
)

const vkSolidityPath = "artifacts/verifier_groth16.sol"

func ExportGroth16VKToSolidity(vk groth16.VerifyingKey) {
	file, err := os.OpenFile(vkSolidityPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(errors.Wrap(err, "failed to open vk file"))
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(errors.Wrap(err, "failed to close vk file"))
		}
	}(file)
	if err := vk.ExportSolidity(file); err != nil {
		panic(errors.Wrap(err, "failed to export vk to solidity"))
	}
}
