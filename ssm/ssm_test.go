package ssm

import "testing"

func TestRun(t *testing.T) {
	result, err := Run("eu-west-1", "sudo id")
	if err != nil {
		t.Errorf("error not expected: err=%s", err)
	}

	if result != "blah" {
		t.Errorf("wrong result: result=%s", result)
	}
}