## GNARK+PRECOMPILE

This project is created to demonstrate how to use [gnark](github.com/Consensys/gnark) 
in composition with forked [geth](https://gitlab.com/distributed_lab/ethereum-experiments/go-geth-precompiled-contracts)
to generate and verify Groth16 proofs using _native_ precompiled contract. 

### Attention

Code base of this repository does not strive to be denominated an epitome of clean code and best engineering practices.
Use it only as a debut on your way to building proof-generating software.

### Trusted setup

For Groth16 zk-SNARK ceremony of generating trusted setup used in this demo project 
should be carefully considered before usage in production.

### Serialization

Considering peculiarities of the ethereum precompiled contracts, and mainly 
[how](https://gitlab.com/distributed_lab/ethereum-experiments/go-geth-precompiled-contracts/-/blob/development/core/vm/contracts.go#L45) 
input is brought in the function, it is needed to have efficient serialization mechanism to pass several arguments in one input. 
What's why here is used the following schema in big-endian format:
* UINT16 of elliptic curve ID 
* UINT32 length of public witness
* Public witness
* UINT32 length of the proof
* Proof 
* Verifier key