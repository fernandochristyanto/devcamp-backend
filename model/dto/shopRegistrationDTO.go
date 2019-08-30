package dto

type ShopRegistrationDTO struct {
	Name        string  `json:"name"`
	Location    string  `json:"location"`
	PhoneNumber string  `json:"phone_number"`
	Email       *string `json:"email"`
	Password    *string `json:"password"`
	UserId      *int32  `json:"user_id"`
}
