package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
)

type Poc struct {
	Deadline    *big.Int `json:"deadline"`
	Nonce       *big.Int `json:"nonce"`
	ScoopNumber *big.Int `json:"scoopNumber"`
}

/**
 * How to un-marshal the Poc struct using the Big.Int rather than the
 * `complexReturn` type.
 */
func (p *Poc) UnmarshalJSON(data []byte) error {
	type Alias Poc
	temp := &struct {
		Deadline    string `json:"deadline"`
		Nonce       string `json:"nonce"`
		ScoopNumber string `json:"scoopNumber"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	deadline, success := big.NewInt(0).SetString(temp.Deadline[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.Deadline))
	}

	nonce, success := big.NewInt(0).SetString(temp.Nonce[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.Nonce))
	}

	scoopNumber, success := big.NewInt(0).SetString(temp.ScoopNumber[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to bigInt", temp.ScoopNumber))
	}

	p.Deadline = deadline
	p.Nonce = nonce
	p.ScoopNumber = scoopNumber

	return nil
}
