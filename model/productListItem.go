package model

type ProductListItem struct {
	ID          int    `json:"id"`
	ShopID      int    `json:"shop_id"`
	CatalogID   int    `json:"catalog_id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
	Charity     bool   `json:"charity"`
	Curated     bool   `json:"curated"`
	ImageUrl    string `json:"image_url"`
	ShopName    string `json:"shop_name"`
}
