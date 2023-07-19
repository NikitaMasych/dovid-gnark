package main

import (
	"main/circuits"
	"main/serialization"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/pkg/errors"
)

func main() {
	var (
		circuit circuits.Poseidon
		ecID    = ecc.BN254
	)

	// ---------------- COMPILE CIRCUIT -------------------
	ccs, err := frontend.Compile(ecID.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(errors.Wrap(err, "failed to compile circuit"))
	}

	// ---------------- SETUP PK AND VK --------------------
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		panic(errors.Wrap(err, "failed to setup pk and vk"))
	}

	// ---------------- DEFINE WITNESS ----------
	hashValue, _ := new(big.Int).SetString("13377623690824916797327209540443066247715962236839283896963055328700043345550", 0)
	assignment := circuits.Poseidon{
		Input: 111,
		Hash:  hashValue,
	}
	wit, err := frontend.NewWitness(&assignment, ecID.ScalarField())
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

	serialization.SerializeAndPrintForGroth16Precompile(proof, vk, wit)
}
