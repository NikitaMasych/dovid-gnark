package utils

import (
	"os"

	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/pkg/errors"
)

const publicWitnessPath = "artifacts/public.json"

func ExportJsonPublicWitness(wit witness.Witness, circuit frontend.Circuit) {
	publicWitness, err := wit.Public()
	if err != nil {
		panic(errors.Wrap(err, "failed to extract public witness"))
	}
	schema, err := frontend.NewSchema(circuit)
	if err != nil {
		panic(errors.Wrap(err, "failed to init new schema"))
	}
	publicWitnessJson, err := publicWitness.ToJSON(schema)
	if err != nil {
		panic(errors.Wrap(err, "failed to transform to json"))
	}
	publicWitnessFile, err := os.OpenFile(publicWitnessPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(errors.Wrap(err, "failed to open public witness file"))
	}
	defer func(publicWitnessFile *os.File) {
		if err := publicWitnessFile.Close(); err != nil {
			panic(errors.Wrap(err, "failed to close public witness file"))
		}
	}(publicWitnessFile)

	if _, err := publicWitnessFile.Write(publicWitnessJson); err != nil {
		panic(errors.Wrap(err, "failed to write public witness json"))
	}
}
