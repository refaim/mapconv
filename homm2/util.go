package homm2

import (
	"../serializers"
)

func (m *Map) Rumors() []*Rumor {
	rumors := make([]*Rumor, m.RumorCount)
	for i := uint8(0); i < m.RumorCount; i++ {
		rumors[i] = &Rumor{}
		reader := serializers.NewByteReader(m.Objects[m.RumorObjectIds[i]].Data)
		rumors[i].Serialize(reader)
	}
	return rumors
}

func (m *Map) Events() []*EventDay {
	events := make([]*EventDay, m.EventCount)
	for i := uint8(0); i < m.EventCount; i++ {
		events[i] = &EventDay{}
		reader := serializers.NewByteReader(m.Objects[m.EventObjectIds[i]].Data)
		events[i].Serialize(reader)
	}
	return events
}

func (m *Map) NameStr() string {
	return extractString(m.Name[:])
}

func (m *Map) DescriptionStr() string {
	return extractString(m.Description[:])
}

func extractString(data []byte) string {
	i := 0
	for data[i] != 0 {
		i++
	}
	return string(data[:i])
}
