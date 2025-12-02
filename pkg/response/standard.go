package response

type ResponseSuccess struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success message"`
	Data    any    `json:"data,omitempty"`
}

type ResponseError struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error,omitempty" example:"Error message"`
}

type HateoasLink struct {
	Self any `json:"self"`
	Next any `json:"next"`
	Prev any `json:"prev"`
	Last any `json:"last"`
}
