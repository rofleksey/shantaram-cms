package api

type IdMessage interface {
	GetId() string
}

func (m *WsOrdersChangedMessage) GetId() string {
	return m.Id
}

func (m *WsMenuChangedMessage) GetId() string {
	return m.Id
}

func (m *WsMessage) GetId() string {
	return m.Id
}
