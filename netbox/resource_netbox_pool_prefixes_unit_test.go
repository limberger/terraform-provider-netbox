package netbox

import (
	// "fmt"
	// "regexp"
	// "strconv"
	"testing"
)

func TestIsLengthValid(t *testing.T) {
	for i := 0; i < netboxMinValidPrefix; i++ {
		if isLengthValid(i) {
			t.Errorf("%d is not a valid length", i)
		}
	}
	for i := netboxMinValidPrefix; i <= netboxMaxValidPrefix; i++ {
		if ! isLengthValid(i) {
			t.Errorf("%d is a valid length", i)
		}
	}
	for i := netboxMaxValidPrefix + 1; i < 35; i++ {
		if isLengthValid(i) {
			t.Errorf("%d is not a valid length", i)
		}
	}
}

func TestIsPoolValid(t *testing.T) {
	good_prefixes := []string{"100.64.0.0/10", "172.16.0.0/12", "10.0.0.0/8"}
	for _, prefix := range good_prefixes {
		if ! isPoolValid(prefix) {
			t.Errorf("%s is a valid pool", prefix)
		}
	}
	// Check some bad ones
	bad_prefixes := []string{"Bad", "172.16.0.0/24", ""}
	for _, prefix := range bad_prefixes {
		if isPoolValid(prefix) {
			t.Errorf("%s is not a valid pool", prefix)
		}
	}
}

func findInSlice(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func TestTagMapToSlice(t *testing.T) {
	tagMap := map[string]interface{} {
		"key1": "value1",
		"key2": "value2",
	}
	expectedSlice := []string{"key1=value1", "key2=value2"}
	generatedSlice := tagMapToSlice(tagMap)
	if len(generatedSlice) != len(expectedSlice) {
		t.Errorf("Bad slice length. Have %d, expected %d", len(generatedSlice), len(expectedSlice))
	}
	for _, expected := range expectedSlice {
		found := findInSlice(generatedSlice, expected)
		if ! found {
			t.Errorf("Did not find %s in generatedSlice", expected)
		}
	}
}

func TestTagSliceToMap(t *testing.T) {
	tagSlice := []string{"mytag=myvalue", "yourtag=yourvalue"}
	expectedMap := map[string]interface{} {
		"mytag": "myvalue",
		"yourtag": "yourvalue",
	}
	generatedMap := tagSliceToMap(tagSlice)
	if len(generatedMap) != len(expectedMap) {
		t.Errorf("Bad map length. Have %d, expected %d", len(generatedMap), len(expectedMap))
	}
	for k, v := range generatedMap {
		found_value, found := expectedMap[k]
		if ! found {
			t.Errorf("%s not found in the generatedMap", k)
		}
		if found_value != v {
			t.Errorf("Key %s has value %s, expected %s", k, found_value, v)
		}
	}
}

func TestIsTagMapValid(t *testing.T) {
	goodMap := map[string]interface{} {
		"name": "testname",
		"key1": "value1",
		"key2": "value2",
		"unique": "testunique",
	}
	if ! isTagMapValid(goodMap) {
		t.Errorf("goodMap reported as bad")
	}

	mapMissingUnique := map[string]interface{} {
		"name": "testname",
		"key1": "value1",
		"key2": "value2",
	}
	if isTagMapValid(mapMissingUnique) {
		t.Errorf("mapMissingUnique reported as good")
	}

	mapMissingName := map[string]interface{} {
		"key1": "value1",
		"key2": "value2",
		"unique": "testunique",
	}
	if isTagMapValid(mapMissingName) {
		t.Errorf("mapMissingName reported as good")
	}

	mapMissingBoth := map[string]interface{} {
		"key1": "value1",
		"key2": "value2",
		"unique": "testunique",
	}
	if isTagMapValid(mapMissingBoth) {
		t.Errorf("mapMissingBoth reported as good")
	}
}
