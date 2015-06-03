package cli

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestNewRequestString(t *testing.T) {
	requestType := "basic"
	requestData := "{\"command\": \"echo hello world\"}"

	request, err := NewRequest(requestType, requestData)
	if err != nil {
		t.Fatal(err)
	}

	if request.Type != requestType || string(request.Data) != requestData {
		t.Fatal("bad request")
	}
}

func TestNewRequestNonExistFile(t *testing.T) {
	requestType := "basic"
	requestData := "./does-not-exist.json"

	request, err := NewRequest(requestType, requestData)
	if err != nil {
		t.Fatal(err)
	}

	if request.Type != requestType || string(request.Data) != requestData {
		t.Fatal("bad request")
	}
}

func TestNewRequestFile(t *testing.T) {
	requestType := "basic"
	expectedData := "content"

	f, err := ioutil.TempFile("/tmp", "quantum-cli")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if _, err := f.Write([]byte(expectedData)); err != nil {
		t.Fatal(f)
	}

	stat, _ := f.Stat()

	request, err := NewRequest(requestType, "/tmp/"+stat.Name())
	if err != nil {
		t.Fatal(err)
	}

	if request.Type != requestType || string(request.Data) != expectedData {
		fmt.Println(string(request.Data))
		fmt.Println(expectedData)
		t.Fatal("bad request")
	}
}
