package mock

import "encoding/json"

// MarshallerStub that will be used for testing
type MarshallerStub struct {
	MarshalCalled   func(obj interface{}) ([]byte, error)
	UnmarshalCalled func(obj interface{}, buff []byte) error
}

// Marshal calls a custom marshall function if defined, otherwise returns json marshal
func (mm *MarshallerStub) Marshal(obj interface{}) ([]byte, error) {
	if mm.MarshalCalled != nil {
		return mm.MarshalCalled(obj)
	}
	return json.Marshal(obj)
}

// Unmarshal calls a custom unmarshall function if defined, otherwise returns json unmarshal
func (mm *MarshallerStub) Unmarshal(obj interface{}, buff []byte) error {
	if mm.UnmarshalCalled != nil {
		return mm.UnmarshalCalled(obj, buff)
	}
	return json.Unmarshal(buff, obj)
}

// IsInterfaceNil returns true if interface is nil, false otherwise
func (mm *MarshallerStub) IsInterfaceNil() bool {
	return mm == nil
}
