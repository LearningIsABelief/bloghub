package request

type LoginUsingPhone struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type LoginUsingPassword struct {
	PhoneOrEmailOrName string `json:"phone_or_email_or_name"`
	Password           string `json:"password"`
}

type PwdReset struct {
	PhoneOrEmail string `json:"phone_or_email"`
	Code         string `json:"code"`
	NewPassword  string `json:"new_password"`
}
