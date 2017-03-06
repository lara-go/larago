package responses

// Response interface.
type Response interface {
	Headers() map[string]string
	SetStatus(int) Response
	Status() int
	ContentType() string
	Body() []byte
}
