package reponse

// Default Error Response
type Error struct {
	Base
}

// Creates a new Error Response
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
