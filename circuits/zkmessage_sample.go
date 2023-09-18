package circuits

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/signature/ecdsa"
	"math/big"
)

func SampleZKMessage() (circuit frontend.Circuit, witness frontend.Circuit) {
	r, _ := new(big.Int).SetString("0x3bca8151a3a1527db0b6afee69c076b0e4c8cab99cf8c5b009f823231f79da6f", 0)
	s, _ := new(big.Int).SetString("0x5ac6bbe39ce0abe42f3bdaae3243717d01cbd7ea9e0865f1f462cc6b8033d222", 0)
	// msg here is like in ethereum signature, keccak256(\x19Ethereum..)
	msg, _ := new(big.Int).SetString("0xf5f6ce66ceaa8a0b97098f19eda6fd65267550f1e7dad91dbbf8d770a151b3fa", 0)
	pubX, _ := new(big.Int).SetString("0xc34e4c872fe7324680f56d6a3e75dc3f6244501ff5bb5dcd12c880180abb3b97", 0)
	pubY, _ := new(big.Int).SetString("0x2ea0298ceb8d8d884e151839967cbfbf8189e06249fd907219cbedb244e82c5d", 0)

	var path [treeLevel]frontend.Variable
	path[0], _ = new(big.Int).SetString("0x74d7d63d7ddada435211e0a92ef2bb2aed9ea61dc92ade0dbe40c352621ec9c", 0)
	for i := 1; i < treeLevel; i++ {
		path[i] = 0 // needed to init path
	}
	rootHash, _ := new(big.Int).SetString("0x1e4ed973c0503efd2602fde8e68549661734162c4ffac402e772c2f31d3a45ef", 0)

	circuit = new(ZKMessage)
	witness = &ZKMessage{
		Sig: ecdsa.Signature[emulated.Secp256k1Fr]{
			R: emulated.ValueOf[emulated.Secp256k1Fr](r),
			S: emulated.ValueOf[emulated.Secp256k1Fr](s),
		},
		Msg: emulated.ValueOf[emulated.Secp256k1Fr](msg),
		Pub: ecdsa.PublicKey[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](pubX),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](pubY),
		},

		HashOrder: 1,
		Path:      path,
		RootHash:  rootHash,
	}

	return
}
