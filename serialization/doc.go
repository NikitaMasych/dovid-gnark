// Package serialization defines a way in which groth16 proofs should be serialized
// in order to be passed to precompiled contract for verification.
// UINT values should be big-endian encoded.
// General scheme:
//
// UINT16 of elliptic curve ID +
//
// UINT32 length of public witness +
//
// Public witness +
//
// UINT32 length of the proof +
//
// Proof +
//
// Verifier key.
package serialization
