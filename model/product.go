package model

type Product struct {
	ID          int    `json:"id"`
	ShopID      int    `json:"shop_id"`
	CatalogID   int    `json:"catalog_id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
	Charity     bool   `json:"charity"`
	Curated     bool   `json:"curated"`
}
