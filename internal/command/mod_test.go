package command

import "testing"

func TestModule(t *testing.T) {
	module, err := Module("../../")
	if err != nil {
		t.Error(err)
		return
	}
	if module != "github.com/syralon/entc-gen-go" {
		t.Fail()
	}
}
