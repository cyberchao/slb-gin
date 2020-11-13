package request

type RegisterAndLoginStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
}
