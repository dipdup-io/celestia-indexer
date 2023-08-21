package handler

import "github.com/btcsuite/btcutil/bech32"

func DecodeAddress(address string) ([]byte, error) {
	_, data, err := bech32.Decode(address)
	return data, err
}
