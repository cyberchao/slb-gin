package request

type LoginStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Iotp     string `json:"iotp"`
	Sms      string `json:"sms"`
}
