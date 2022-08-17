package request

type Update struct {
	Code            string `json:"code"`
	NewPhoneOrEmail string `json:"new_phone_or_email"`
}
