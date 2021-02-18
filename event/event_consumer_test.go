package event

import (
	"testing"
)

func TestAddCallback(t *testing.T) {
	err := AddCallback("COMMAND_EVENT", func(message ConsumerMessage) error {
		return nil
	})

	if  err != nil {
		t.Error("Fail to test, Add Callback return error")
	}

	err = AddCallback("COMMAND_EVENT", func(message ConsumerMessage) error {
		return nil
	})

	if err == nil {
		t.Error("Fail to test, Add callback should return error")
	}
}

func TestRemoveCallback(t *testing.T) {
	err := RemoveCallback("COMMAND_EVENT")

	if err != nil {
		t.Error("Fail to test, Remove callback shouldn't return error")
	}

	err = RemoveCallback("COMMAND_EVENT2")

	if err == nil {
		t.Error("Fail to test, Remove callback should return error")
	}
}

func TestGenerateConsumerMessage(t *testing.T) {

}