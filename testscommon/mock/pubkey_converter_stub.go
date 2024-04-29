package mock

// PubKeyConverterStub that will be used for testing
type PubKeyConverterStub struct {
	DecodeCalled func(humanReadable string) ([]byte, error)
	EncodeCalled func(pkBytes []byte) string
}

// Len returns zero
func (pkcs *PubKeyConverterStub) Len() int {
	return 0
}

// Decode calls a custom decode function if defined, otherwise returns an empty byte slice, nil err
func (pkcs *PubKeyConverterStub) Decode(humanReadable string) ([]byte, error) {
	if pkcs.DecodeCalled != nil {
		return pkcs.DecodeCalled(humanReadable)
	}

	return make([]byte, 0), nil
}

// Encode calls a custom encode function if defined, otherwise returns "moa1"+input, nil err
func (pkcs *PubKeyConverterStub) Encode(pkBytes []byte) string {
	if pkcs.EncodeCalled != nil {
		return pkcs.EncodeCalled(pkBytes)
	}

	return "moa1" + string(pkBytes)
}

// IsInterfaceNil returns true if interface is nil, false otherwise
func (pkcs *PubKeyConverterStub) IsInterfaceNil() bool {
	return pkcs == nil
}
