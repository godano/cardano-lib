package chain

import cardano "github.com/godano/cardano-lib"

type Transaction struct {
	_       struct{} `cbor:",toarray"`
	Body    TransactionBody
	Witness TransactionWitness
}

// TransactionBody
type TransactionBody struct {
	Inputs  []TransactionInput  `cbor:"0,keyasint"`
	Outputs []TransactionOutput `cbor:"1,keyasint"`
	Fee     cardano.Coin        `cbor:"2,keyasint"`
	TTL     uint64              `cbor:"3,keyasint"`
}

// TransactionInput
type TransactionInput struct {
	_     struct{} `cbor:",toarray"`
	ID    []byte
	Index uint64
}

// TransactionOutput specifies an amount of coins and an address to which this amount shall be sent.
type TransactionOutput struct {
	_       struct{} `cbor:",toarray"`
	Address []byte
	Amount  cardano.Coin
}

type TransactionWitness struct {
}

type TransactionMetadata struct {
}

//
type TransactionBuilder interface {

	// Build
	Build() (*Transaction, error)
}

type simpleBuilder struct {
}

func NewTransactionBuilder() *TransactionBuilder {
	return nil
}

func (builder *simpleBuilder) Build() (*Transaction, error) {

	return nil, nil
}
