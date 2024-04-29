package mock

// ShardCoordinatorMock -
type ShardCoordinatorMock struct {
	SelfID uint32
}

// ComputeId returns 0
func (scm *ShardCoordinatorMock) ComputeId([]byte) uint32 {
	return 0
}

// SelfId returns SelfID member
func (scm *ShardCoordinatorMock) SelfId() uint32 {
	return scm.SelfID
}

// IsInterfaceNil returns true if interface is nil, false otherwise
func (scm *ShardCoordinatorMock) IsInterfaceNil() bool {
	return scm == nil
}
