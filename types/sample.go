package types

type SamplePayload struct {
	Text string `json:"text" validate:"required|TextValidator" message:"text is required and should be abhinav"`
}

func (f SamplePayload) TextValidator(val string) bool {
	// get comparison values from constanst
	return (val == "abhinav")
}
