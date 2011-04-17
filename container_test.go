package avrogo

import (
	"os"
	"testing"
)

func TestWeatherContainer(t *testing.T) {
	reader, err := os.Open("weather.avro", os.O_RDONLY, 0444)
	if err != nil {
		t.Fatalf("Unable to open 'weather.avro'")
	}
	_, iter, err := ReadContainer(reader)

	checkNextTemp := func(expected int32) {
		extracted := iter.Next().(map[string]interface{})["temp"].(int32)
		if extracted != expected {
			t.Fatalf("Expected %v but got %v", expected, extracted)
		}
	}

	checkNextTemp(0)
	checkNextTemp(22)
	checkNextTemp(-11)
	checkNextTemp(111)
	checkNextTemp(78)
	if hasNext, err := iter.HasNext(); hasNext {
		t.Fatal("There are only 5 values, but hasNext is still true")
	} else if err.(os.Error).String() != "EOF" {
		t.Fatal("Expecting err to be EOF")
	}
}
