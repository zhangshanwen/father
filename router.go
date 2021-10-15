package father

type (
	Router struct {
		Method   string
		Path     string
		Handlers []HandlerFunc
	}
)
