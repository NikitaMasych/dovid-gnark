package circuits

import "github.com/consensys/gnark/frontend"

// Circuit0 defines a simple circuit
// x**3 + x + 5 == y
type Circuit0 struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
}

// Circuit1 defines a simple circuit
// x**3 + x + 5 == y + z
type Circuit1 struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
	Z frontend.Variable `gnark:",public"`
}

// Circuit2 defines a simple circuit
// x**3 + x + 5 == (y + z) * c
type Circuit2 struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
	Z frontend.Variable `gnark:",public"`
	C frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
// x**3 + x + 5 == y
func (circuit *Circuit0) Define(api frontend.API) error {
	x3 := api.Mul(circuit.X, circuit.X, circuit.X)
	api.AssertIsEqual(circuit.Y, api.Add(x3, circuit.X, 5))
	return nil
}

// Define declares the circuit constraints
// x**3 + x + 5 == y + z
func (circuit *Circuit1) Define(api frontend.API) error {
	x3 := api.Mul(circuit.X, circuit.X, circuit.X)
	api.AssertIsEqual(api.Add(circuit.Y, circuit.Z), api.Add(x3, circuit.X, 5))
	return nil
}

// Define declares the circuit constraints
// x**3 + x + 5 == (y + z) * c
func (circuit *Circuit2) Define(api frontend.API) error {
	x3 := api.Mul(circuit.X, circuit.X, circuit.X)
	temp := api.Add(circuit.Y, circuit.Z)
	api.AssertIsEqual(api.Mul(temp, circuit.C), api.Add(x3, circuit.X, 5))
	return nil
}
