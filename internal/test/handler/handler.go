package handler

type TestService interface {
}

type Handler struct {
	testService TestService
}

func NewHandler(testService TestService) *Handler {
	return &Handler{testService: testService}
}
