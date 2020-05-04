package Bdiamond

import (
	"blockbook/bchain/coins/btc"

	"github.com/martinboehm/btcd/wire"
	"github.com/martinboehm/btcutil/chaincfg"
)

// magic numbers
const (
	MainnetMagic wire.BitcoinNet = 0xf9beb4d9
	TestnetMagic wire.BitcoinNet = 0x0b110907
)

// chain parameters
var (
	MainNetParams chaincfg.Params
	TestNetParams chaincfg.Params
)

func init() {
	MainNetParams = chaincfg.MainNetParams
	MainNetParams.Net = MainnetMagic

	// Address encoding magics
	MainNetParams.PubKeyHashAddrID = []byte{145}
	MainNetParams.ScriptHashAddrID = []byte{23}

	TestNetParams = chaincfg.TestNet3Params
	TestNetParams.Net = TestnetMagic

	// Address encoding magics
	TestNetParams.PubKeyHashAddrID = []byte{145}
	TestNetParams.ScriptHashAddrID = []byte{23}
}

// BdiamondParser handle
type BdiamondParser struct {
	*btc.BitcoinParser
}

// NewBdiamondParser returns new BdiamondParser instance
func NewBdiamondParser(params *chaincfg.Params, c *btc.Configuration) *BdiamondParser {
	return &BdiamondParser{BitcoinParser: btc.NewBitcoinParser(params, c)}
}

// GetChainParams contains network parameters for the main Bdiamond network,
// the regression test Bdiamond network, the test Bdiamond network and
// the simulation test Bdiamond network, in this order
func GetChainParams(chain string) *chaincfg.Params {
	if !chaincfg.IsRegistered(&MainNetParams) {
		err := chaincfg.Register(&MainNetParams)
		if err == nil {
			err = chaincfg.Register(&TestNetParams)
		}
		if err != nil {
			panic(err)
		}
	}
	switch chain {
	case "test":
		return &TestNetParams
	default:
		return &MainNetParams
	}
}
