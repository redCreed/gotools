package queue

type Message struct {
	ID         string
	Stream     string
	Values     map[string]interface{}
	ErrorCount int
}

func (m *Message) GetID() string {
	return m.ID
}

func (m *Message) GetStream() string {
	return m.Stream
}

func (m *Message) GetValues() map[string]interface{} {
	return m.Values
}

func (m *Message) SetID(id string) {
	m.ID = id
}

func (m *Message) SetStream(stream string) {
	m.Stream = stream
}

func (m *Message) SetValues(values map[string]interface{}) {
	m.Values = values
}

func (m *Message) GetPrefix() (prefix string) {
	//	if m.Values == nil {
	//		return
	//	}
	//	v, _ := m.Values[storage.PrefixKey]
	//	prefix, _ = v.(string)
	return
}

func (m *Message) SetPrefix(prefix string) {
	//	if m.Values == nil {
	//		m.Values = make(map[string]interface{})
	//	}
	//	m.Values[storage.PrefixKey] = prefix
	return
}

func (m *Message) SetErrorCount(count int) {
	m.ErrorCount = count
}

func (m *Message) GetErrorCount() int {
	return m.ErrorCount
}
