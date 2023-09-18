package circuits

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/hash/sha3"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/std/signature/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/vocdoni/gnark-crypto-primitives/poseidon"
	"main/merkle"
	"math/big"
)

const treeLevel = 32

type ZKMessage struct {
	Sig ecdsa.Signature[emulated.Secp256k1Fr]
	Msg emulated.Element[emulated.Secp256k1Fr] `gnark:"msg,public"`
	Pub ecdsa.PublicKey[emulated.Secp256k1Fp, emulated.Secp256k1Fr]

	HashOrder frontend.Variable
	Path      [treeLevel]frontend.Variable
	RootHash  frontend.Variable `gnark:"root,public"`
}

func (c *ZKMessage) Define(api frontend.API) error {
	c.Pub.Verify(api, sw_emulated.GetSecp256k1Params(), &c.Msg, &c.Sig)

	pubKeyBytes, err := pubKeyToBytes(api, &c.Pub)
	if err != nil {
		return err
	}
	keccak, err := sha3.NewLegacyKeccak256(api)
	if err != nil {
		return err
	}
	keccak.Write(pubKeyBytes)
	digest := keccak.Sum()
	addressBytes := digest[keccak.Size()-common.AddressLength:]
	leaf := bytesToAddress(api, addressBytes)

	mp := merkle.MerkleProof{
		RootHash: c.RootHash,
		Path:     append([]frontend.Variable{leaf}, c.Path[:]...),
	}
	h := poseidon.NewPoseidon(api)
	mp.VerifyProof(api, &h, c.HashOrder)

	return nil
}

func bytesToAddress(api frontend.API, array []uints.U8) frontend.Variable {
	v := make([]frontend.Variable, common.AddressLength)
	for i := range v {
		bitsShift := (common.AddressLength-i)*8 - 8
		multiplier := new(big.Int).Lsh(big.NewInt(1), uint(bitsShift))
		v[i] = api.Mul(array[i].Val, multiplier)
	}
	return api.Add(v[0], v[1], v[2:]...)
}

func pubKeyToBytes(api frontend.API, pubKey *ecdsa.PublicKey[emulated.Secp256k1Fp, emulated.Secp256k1Fr]) ([]uints.U8, error) {
	xLimbs := pubKey.X.Limbs
	yLimbs := pubKey.Y.Limbs

	u64api, err := uints.New[uints.U64](api)
	if err != nil {
		return nil, err
	}

	result := limbsToBytes(u64api, xLimbs)
	return append(result, limbsToBytes(u64api, yLimbs)...), nil
}

func limbsToBytes(u64api *uints.BinaryField[uints.U64], limbs []frontend.Variable) []uints.U8 {
	result := make([]uints.U8, 0, len(limbs)*8)
	for i := range limbs {
		u64 := u64api.ValueOf(limbs[len(limbs)-1-i])
		result = append(result, u64api.UnpackMSB(u64)...)
	}
	return result
}
