package dto

type GarageSaleProductInsertDTO struct {
	ShopID      int    `json:"shop_id"`
	CatalogID   int    `json:"catalog_id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
	Charity     bool   `json:"charity"`
}
