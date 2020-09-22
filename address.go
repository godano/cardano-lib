package cardano_lib

import "errors"

type AddressType uint

const (
	Payment AddressType = iota + 1
	Pointer
	Reward
	Enterprise
	Byron
)

var (
	// UnknownAddressTypeError is thrown, if an unknown AddressType is encountered.
	UnknownAddressTypeError = errors.New("unknown address type")
)

// Address is an interface type for any type of address in the Cardano
// blockchain, including shelley base address, pointer, enterprise
// reward and byron base address. This interface is designed to be
// generic enough such that new addresses can be added later.
type Address interface {
	// Type returns the type of the address
	Type() AddressType
	// NetworkID returns the ID of the blockchain network
	NetworkID() uint8
}

type addressImpl struct {
	_type     AddressType
	networkID uint8
}

// Type returns the type of the address
func (address addressImpl) Type() AddressType {
	return address._type
}

// NetworkID returns the ID of the blockchain network
func (address addressImpl) NetworkID() uint8 {
	return address.networkID
}

// PaymentAddress
type PaymentAddress struct {
	addressImpl
}

// MarshalCBOR marshals a PaymentAddress into its corresponding CBOR representation. A payment address
// in the Shelley era has the following format: [ 8 bit header | payload ]
//
// Header:
//   bits 3-0: network id
//   bits   4: payment cred is keyhash/scripthash
//   bits   5: pointer/enterprise [for base: stake cred is keyhash/scripthash]
//   bits   6: base/other
//   bits   7: 0
//
// Payload:
func (address *PaymentAddress) MarshalCBOR() (data []byte, err error) {
	return nil, nil
}

func (address *PaymentAddress) UnmarshalCBOR([]byte) error {
	return nil
}
