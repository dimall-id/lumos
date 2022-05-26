package errs

import "fmt"

type ErrDecodeJson struct {
	dueTo error
}

func (e ErrDecodeJson) Error() string {
	return fmt.Sprintf("fail to decode to json string due to '%s'", e.dueTo)
}

func DecodeJsonErr (dueTo error) ErrDecodeJson {
	return ErrDecodeJson{
		dueTo: dueTo,
	}
}

type ErrGenerateOutbox struct {
	dueTo error
}

func (e ErrGenerateOutbox) Error() string {
	return fmt.Sprintf("fail to generate outbox message due to '%s'", e.dueTo)
}

func GenerateOutboxErr (dueTo error) ErrGenerateOutbox {
	return ErrGenerateOutbox{
		dueTo: dueTo,
	}
}
