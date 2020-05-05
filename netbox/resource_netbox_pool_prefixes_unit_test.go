package netbox

import (
	"fmt"
	"testing"
)

// Test the envs that match the vrfs exactly
var goodEnvs = [4](string){"dev", "stage", "test", "prod"}
var goodResources = [5](string){"core", "depot", "edge", "servicedelivery", "vpn"}

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

func TestEnvToVrf(t *testing.T) {
	for _, env := range goodEnvs {
		vrf := envToVrf(env)
		if vrf != env {
			t.Errorf("envToVrf returned %s expected %s", vrf, env)
		}
	}

	// Test with mixed case
	mixedEnvs := [4](string){"Dev", "sTage", "Test", "prOd"}
	for i, env := range mixedEnvs {
		vrf := envToVrf(env)
		if vrf != goodEnvs[i] {
			t.Errorf("envToVrf returned %s expected %s", vrf, goodEnvs[i])
		}
	}

	// Test for predev
	vrf := envToVrf("non-standard")
	if vrf != "pre-dev" {
		t.Errorf("envToVrf returned %s expected pre-dev", vrf)
	}

}

func TestGetVrf(t *testing.T) {
	// This really just tests exceptions for servicedelivery.
	// All other relevant tests are in TestEnvToVrf.
	// Test for servicedelivery
	for _, env := range goodEnvs {
		vrf := getVrf(env, "servicedelivery")
		if vrf != "servicedelivery" {
			t.Errorf("getVrf returned %s expected servicedelivery", vrf)
		}
	}

	// A non-standard environment--even for servicedelivery, should use pre-dev
	vrf := getVrf("non-standard", "servicedelivery")
	if vrf != "pre-dev" {
		t.Errorf("getVrf returned %s expected pre-dev", vrf)
	}
}

func verifySupernets(rType string, found [](string), expected [](string), t *testing.T) {
	if len(found) != len(expected) {
		t.Errorf("Found %d supernets for %s. Expected %d", len(found), rType, len(expected))
	}
	for _, f := range found {
		matched := false
		for _, e := range expected {
			if e == f {
				matched = true
				break
			}
		}
		if ! matched {
			t.Errorf("Supernet %s for '%s' not in expected set %v", f, rType, expected)
		}
	}
}

func TestResourceTypeToSupernets(t *testing.T) {

	found := resourceTypeToSupernets("core")
	verifySupernets("core", found, [](string){"100.64.0.0/10"}, t)

	found = resourceTypeToSupernets("depot")
	verifySupernets("depot", found, [](string){"10.224.0.0/16", "10.225.0.0/16"}, t)

	found = resourceTypeToSupernets("edge")
	verifySupernets("edge", found, [](string){"10.226.0.0/16", "10.227.0.0/16"}, t)

	found = resourceTypeToSupernets("servicedelivery")
	verifySupernets("servicedelivery", found, [](string){"10.228.0.0/16"}, t)

	found = resourceTypeToSupernets("vpn")
	verifySupernets("vpn", found, [](string){"172.16.0.0/12"}, t)

	found = resourceTypeToSupernets("bad")
	if len(found) != 0 {
		t.Errorf("Supernet %v for 'bad' was not expected", found)
	}
}

func TestAreResourceAndEnvValid(t *testing.T) {
	allResources := []string{"servicedelivery", "core", "edge", "depot", "vpn"}
	allGoodEnvs := []string{"pre-dev", "dev", "stage", "test", "prod"}

	for _, r := range allResources {
		for _, e := range allGoodEnvs {
			if ! areResourceAndEnvValid(r, e) {
				t.Errorf(fmt.Sprintf("%s resource should be OK with %s environment", r, e))
			}
		}
	}
	for _, r := range allResources {
		// No other resource type should work with servicedelivery env
		valid := areResourceAndEnvValid(r, "servicedelivery")

		if r == "servicedelivery" {
			if ! valid {
				t.Errorf("servicedelivery resource should be OK with servicedelivery environment")
			}
		} else {
			if valid {
				t.Errorf(fmt.Sprintf("%s resource should not be OK with servicedelivery environment", r))
			}
		}
	}
}
