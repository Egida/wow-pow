package protocol

type Quote struct {
	Quote string
}

func NewQuote(quote string) *Quote {
	return &Quote{
		Quote: quote,
	}
}

func (m *Quote) Size() uint64 {
	return uint64(len(m.Quote))
}

func (m *Quote) MessageType() MessageType {
	return MessageTypeQuote
}

func (m *Quote) Serialize(buf []byte) {
	copy(buf, m.Quote)
}

func (m *Quote) Deserialize(buf []byte) {
	m.Quote = string(buf)
}
