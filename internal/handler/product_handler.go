package handler

import (
	"encoding/json"
	"fmt"
	internal "github.com/fernandochristyanto/devcamp-backend/internal"
	"github.com/fernandochristyanto/devcamp-backend/model"
	"github.com/fernandochristyanto/devcamp-backend/model/dto"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GetGarageSales(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	garageSalesQuery := fmt.Sprintf("SELECT id, shop_id, catalog_id, name, price, description, stock, charity, curated FROM products where catalog_id = 1 ORDER BY curated = true desc")

	rows, err := h.DB.Query(garageSalesQuery)

	var products []model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ID,
			&product.ShopID,
			&product.CatalogID,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Stock,
			&product.Charity,
			&product.Curated,
		)

		if err != nil {
			log.Println("Error scanning list products")
			continue
		}
		products = append(products, product)
	}

	bytes, err := json.Marshal(products)
	if err != nil {
		log.Println(err)
		return
	}
	internal.RenderJSON(w, bytes, http.StatusOK)
}

func (h *Handler) GetProductDetail(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	shopId := param.ByName("id")
	selectShopQuery := fmt.Sprintf("SELECT id, shop_id, catalog_id, name, price, description, stock, charity, curated FROM products where id = %s", shopId)

	rows, err := h.DB.Query(selectShopQuery)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ID,
			&product.ShopID,
			&product.CatalogID,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Stock,
			&product.Charity,
			&product.Curated,
		)

		if err != nil {
			log.Println("Error scanning list products")
			continue
		}
		bytes, err := json.Marshal(product)
		internal.RenderJSON(w, bytes, http.StatusOK)
		return
	}
}

func (h *Handler) NewGarageSaleProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		internal.RenderJSON(w, []byte(`"message": "failed to read body"`), http.StatusBadRequest)
	}

	var newProductDTO dto.GarageSaleProductInsertDTO
	err = json.Unmarshal(body, &newProductDTO)
	if err != nil {
		internal.RenderJSON(w, []byte(`"message": "failed parsing product"`), http.StatusBadRequest)
	}

	// Insert new user (seller)
	insertProductQuery := fmt.Sprintf("INSERT INTO products(shop_id, catalog_id, name, price, description, stock, charity) VALUES(%d, %d, '%s', %d,'%s', %d, %s); SELECT max(id) from shops",
		newProductDTO.ShopID,
		newProductDTO.CatalogID,
		newProductDTO.Name,
		newProductDTO.Price,
		newProductDTO.Description,
		newProductDTO.Stock,
		strconv.FormatBool(newProductDTO.Charity),
	)

	_, err = h.DB.Exec(insertProductQuery)
	if err != nil {
		log.Println(err)
		internal.RenderJSON(w, []byte(`"message": "failed insert product"`), http.StatusBadRequest)
	}

	internal.RenderJSON(w, []byte(`"message": "Success"`), http.StatusOK)
}
