package event

import "testing"

func TestAddCallback(t *testing.T) {
	AddCallback("COMMAND_EVENT", func(topic string, data string) error {
		return nil
	})

	if _, oke := callbacks["COMMAND_EVENT"]; !oke {

	}
}

func TestRemoveCallback(t *testing.T) {

}

