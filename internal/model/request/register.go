package request

type Register struct {
	PhoneOrEmail string `json:"phone_or_email"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	Code         string `json:"code"`
	Age          int    `json:"age"`
}
