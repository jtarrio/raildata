package errors

var (
	BadCredentialsError     error = &badCredentialsError{}     // invalid username or password.
	MissingCredentialsError error = &missingCredentialsError{} // token not present or malformed.
	InvalidTokenError       error = &invalidTokenError{}       // invalid token.
)

// NewRailDataError reports an error produced by the RailData API.
func NewRailDataError(message string) error {
	return &RailDataError{message: message}
}

// RailDataError contains an error produced by the RailData API.
type RailDataError struct {
	message string
}

func (e *RailDataError) Error() string {
	return e.message
}

type badCredentialsError struct{}

func (e *badCredentialsError) Error() string {
	return "invalid username or password in request"
}

type missingCredentialsError struct{}

func (e *missingCredentialsError) Error() string {
	return "missing or malformed credentials in request"
}

type invalidTokenError struct{}

func (e *invalidTokenError) Error() string {
	return "invalid token in request"
}
