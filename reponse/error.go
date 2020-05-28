package reponse

type Error struct {
	Base
}

func NewError(code int, message interface{}) Error {
	data := map[string]interface{}{
		"error":   true,
		"message": message,
	}
	return Error{
		Base{
			code,
			data,
		},
	}
}
