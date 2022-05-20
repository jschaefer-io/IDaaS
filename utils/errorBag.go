package utils

type ErrorBag struct {
	hasErrors bool
	errors    map[string]string
}

func NewErrorBag() *ErrorBag {
	return &ErrorBag{
		hasErrors: false,
		errors:    make(map[string]string),
	}
}

func (b *ErrorBag) AddError(key string, message string) {
	b.hasErrors = true
	b.errors[key] = message
}

func (b *ErrorBag) Empty() bool {
	return !b.hasErrors
}

func (b ErrorBag) Errors() map[string]string {
	list := make(map[string]string)
	for k, v := range b.errors {
		list[k] = v
	}
	return list
}
