package mock

import "io"

type WSConnStub struct {
	io.Closer
	WriteMessageCalled func(messageType int, data []byte) error
	ReadMessageCalled  func() (messageType int, p []byte, err error)
	CloseCalled        func() error
}

func (wsc *WSConnStub) ReadMessage() (messageType int, p []byte, err error) {
	if wsc.ReadMessageCalled != nil {
		return wsc.ReadMessageCalled()
	}
	return 0, nil, nil
}

func (wsc *WSConnStub) WriteMessage(messageType int, data []byte) error {
	if wsc.WriteMessageCalled != nil {
		return wsc.WriteMessageCalled(messageType, data)
	}
	return nil
}

func (wsc *WSConnStub) Close() error {
	if wsc.CloseCalled != nil {
		return wsc.CloseCalled()
	}
	return nil
}
