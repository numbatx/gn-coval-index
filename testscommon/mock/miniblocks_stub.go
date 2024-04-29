package mock

import (
	"github.com/numbatx/gn-coval-index/schema"
	"github.com/numbatx/gn-core/data"
)

// MiniBlockHandlerStub that will be used for testing
type MiniBlockHandlerStub struct {
	ProcessMiniBlockCalled func(header data.HeaderHandler, body data.BodyHandler) ([]*schema.MiniBlock, error)
}

// ProcessMiniBlocks calls a custom mini blocks process function if defined, otherwise returns nil, nil
func (mbhs *MiniBlockHandlerStub) ProcessMiniBlocks(header data.HeaderHandler, body data.BodyHandler) ([]*schema.MiniBlock, error) {
	if mbhs.ProcessMiniBlockCalled != nil {
		return mbhs.ProcessMiniBlockCalled(header, body)
	}

	return nil, nil
}
