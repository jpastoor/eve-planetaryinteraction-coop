package main

import (
	"testing"
	"io/ioutil"
	"reflect"
)

func TestParse(t *testing.T) {
	tp := NewTransactionParser(&TypeFetcherMock{
		cacheByName: map[string]*Type{
			"Oxygen": {TypeID: 1, TypeName: "Oxygen", Volume: 1},
			"Nano-Factory": {TypeID: 2, TypeName: "Nano-Factory", Volume: 2},
			"Reactive Metals": {TypeID: 3, TypeName: "Reactive Metals", Volume: 3},
		},
	})

	file, _ := ioutil.ReadFile("./examples/example_log.log")
	output, errs := tp.Parse(string(file))

	if errs != nil {
		t.Fatalf("Did not expect error, got %s", errs)
	}

	// We compare the hashes since thats also an easy way to ensure the internals are good without any pointer mumbojumbo
	expected := []string{
		"e14cb9eee71975272b4af0060dd17d3b", "248e5484dc72a6527630a2af490ced6b", "e859afbb544d2bcc215fea09e72e7878",
		"339a794f88024bea038ac39074187d10", "c0959f81b8a543173c9e5e373836f0d7", "0ad0631b4c4cf304afa903ac17f7c1be",
		"7ea928959682a8c4a97461f8ece9e453",
	}

	var outputHashes []string
	for _, ts := range output {
		outputHashes = append(outputHashes, ts.Id)
	}

	if !reflect.DeepEqual(expected, outputHashes) {
		t.Errorf("Expected %v, but got %v", expected, outputHashes)
	}
}
