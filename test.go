package main

import (
	"./homm2"
	"./serializers"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	data, err := ioutil.ReadFile("maps/Map_0001.MP2")
	if err != nil {
		log.Fatal(err)
	}

	reader := serializers.NewByteReader(data)
	m := &homm2.Map{}
	m.Serialize(reader)

	// fmt.Printf("%# v\n", m.Map)

	fmt.Printf("Bytes left: %d\n", len(data)-reader.Pos())

	writer := serializers.NewByteWriter()
	m.Serialize(writer)
	savedData := writer.Data()

	if err := ioutil.WriteFile("maps/TestMap.MP2", savedData, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original length = %d, saved length = %d\n", reader.Pos(), len(savedData))

	f := len(data)
	if len(savedData) < f {
		f = len(savedData)
	}

	for i := 0; i < f; i++ {
		if data[i] != savedData[i] {
			fmt.Printf("Diff at pos %d\n", i)
			fmt.Println(hex.Dump(data[i-10 : i+100]))
			fmt.Println(hex.Dump(savedData[i-10 : i+100]))
			break
		}
	}

	fmt.Println(m.RumorCount)
	for _, rumor := range m.Rumors() {
		fmt.Printf("%#v\n", rumor)
	}

	for _, event := range m.Events() {
		fmt.Printf("%#v\n", event)
	}

	fmt.Printf("Name = %s, Description = %s\n", m.NameStr(), m.DescriptionStr())
}
