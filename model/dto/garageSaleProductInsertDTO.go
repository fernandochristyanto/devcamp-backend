package dto

type GarageSaleProductInsertDTO struct {
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
	Charity     bool   `json:"charity"`
}
