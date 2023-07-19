package circuits

import (
	"main/poseidon"

	"github.com/consensys/gnark/frontend"
)

type Poseidon struct {
	Input frontend.Variable `gnark:"input"`
	Hash  frontend.Variable `gnark:",public"`
}

func (circuit *Poseidon) Define(api frontend.API) error {
	poseidonHash := poseidon.NewPoseidon1(api)
	poseidonHash.Write(circuit.Input)
	api.AssertIsEqual(circuit.Hash, poseidonHash.Sum())
	return nil
}
