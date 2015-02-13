package homm2

import (
	"../serializers"
	"fmt"
)

type MapObject struct {
	X      uint8
	Y      uint8
	Object interface{}
}

type Jail struct {
	Hero
}

type RandomCastle struct {
	Castle
}

type BottleObject struct {
	Info
}

type SignObject struct {
	Info
}

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

func (m *Map) MapObjects() (result []*MapObject) {
	fmt.Println("Objects Len =", len(m.Objects))
	for i, tile := range m.Tiles {
		x := uint8(i % int(m.Width))
		y := uint8(i / int(m.Width))
		orders := int(tile.Quantity1) | (int(tile.Quantity2) << 8)
		objectIndex := orders / 8
		if tile.GeneralObject > 0 && orders%8 == 0 && objectIndex > 0 {
			reader := serializers.NewByteReader(m.Objects[objectIndex].Data)

			var obj interface {
				Serialize(Serializer)
			}

			switch tile.GeneralObject {
			case ObjectTypeRandomTown:
				obj = &RandomCastle{}
			case ObjectTypeRandomCastle:
				obj = &RandomCastle{}
			case ObjectTypeCastle:
				obj = &Castle{}
			case ObjectTypeHeroes:
				obj = &Hero{}
			case ObjectTypeSign:
				obj = &SignObject{}
			case ObjectTypeBottle:
				obj = &BottleObject{}
			case ObjectTypeEvent:
				obj = &EventCoord{}
			case ObjectTypeSphinx:
				obj = &Riddle{}
			case ObjectTypeJail:
				obj = &Jail{}
			}

			if obj != nil {
				obj.Serialize(reader)
				result = append(result, &MapObject{X: x, Y: y, Object: obj})
			}
		}
	}
	return
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
