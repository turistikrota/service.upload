package req

type Request interface{}

type request struct{}

func New() Request {
	return &request{}
}
