package mock

// HasherMock that will be used for testing
type HasherMock struct {
}

// Compute outputs a constant dummy hash
func (hs *HasherMock) Compute(string) []byte {
	return []byte("ok")
}

// Size returns a dummy size
func (hs *HasherMock) Size() int {
	return 123
}

// IsInterfaceNil returns true if interface is nil, false otherwise
func (hs *HasherMock) IsInterfaceNil() bool {
	return hs == nil
}
