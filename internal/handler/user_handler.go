package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fernandochristyanto/devcamp-backend/model"

	"github.com/fernandochristyanto/devcamp-backend/internal"
	"github.com/fernandochristyanto/devcamp-backend/model/dto"
	"github.com/julienschmidt/httprouter"
)

func (h *Handler) SellerRegistration(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		internal.RenderJSON(w, []byte(`"message": "failed to read body"`), http.StatusBadRequest)
	}

	var sellerRegistrationDTO dto.ShopRegistrationDTO
	err = json.Unmarshal(body, &sellerRegistrationDTO)
	if err != nil {
		internal.RenderJSON(w, []byte(`"message": "failed parsing user"`), http.StatusBadRequest)
	}

	// Insert new user (seller)
	if sellerRegistrationDTO.UserId == nil || *sellerRegistrationDTO.UserId == 0 {
		insertedUserId := createUser(h.DB, *sellerRegistrationDTO.Email, *sellerRegistrationDTO.Password, "seller", sellerRegistrationDTO.PhoneNumber)
		sellerRegistrationDTO.UserId = &insertedUserId
	}

	insertShopQuery := fmt.Sprintf("INSERT INTO shops(user_id, name, location) VALUES(%d, '%s', '%s'); SELECT max(id) from shops",
		*sellerRegistrationDTO.UserId,
		sellerRegistrationDTO.Name,
		sellerRegistrationDTO.Location,
	)
	rows, err := h.DB.Query(insertShopQuery)
	var insertedShopId int
	for rows.Next() {
		err := rows.Scan(
			&insertedShopId,
		)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	// Insert garage sale catalog
	insertGarageSaleCatalogQuery := fmt.Sprintf("INSERT INTO catalogs(shop_id, name) values(%d, 'Garage Sale')", insertedShopId)
	_, err = h.DB.Exec(insertGarageSaleCatalogQuery)

	user := getUserByID(h.DB, *sellerRegistrationDTO.UserId)

	internal.RenderJSON(w, []byte(fmt.Sprintf(`
	{
		"message": "Success",
		"userId" : %d,
		"email" : "%s",
		"password" : "%s"
	}`, *sellerRegistrationDTO.UserId, user.Email, user.Password)), http.StatusOK)
}

func createUser(db *sql.DB, email string, password string, role string, phone string) int32 {
	query := fmt.Sprintf("INSERT INTO users (email, password, role, phone) VALUES('%s', '%s', '%s', '%s'); SELECT max(id) from users;",
		email,
		password,
		role,
		phone,
	)

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return 0
	}

	var insertedUserId int32
	for rows.Next() {
		err := rows.Scan(
			&insertedUserId,
		)
		if err != nil {
			log.Println(err)
			continue
		}
		return insertedUserId
	}

	return 0
}

func getUserByID(db *sql.DB, id int32) model.User {
	query := fmt.Sprintf("Select id,email,password from users where id = %d", id)
	rows, err := db.Query(query)
	if err != nil {
		println("gagal")
	}
	defer rows.Close()
	var user model.User
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.Password)

		if err != nil {
			println(err)
			continue
		}

	}

	return user

}

func (h *Handler) GetProductsByUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userId := params.ByName("id")
	query := fmt.Sprintf(`SELECT d.id, d.shop_id, d.catalog_id,
	d.name, d.price, d.description, d.stock,
	d.charity, d.curated, d.image_url from 
	users a join shops b on a.id = b.user_id 
	join catalogs c on c.shop_id = b.id  
	join products d on d.catalog_id = c.id 
	where a.id = %s`, userId)

	rows, err := h.DB.Query(query)
	if err != nil {
		log.Println("Error selecting products by user")
	}

	var products []model.ProductListItem
	for rows.Next() {
		var product model.ProductListItem
		err = rows.Scan(
			&product.ID,
			&product.ShopID,
			&product.CatalogID,
			&product.Name,
			&product.Price,
			&product.Description,
			&product.Stock,
			&product.Charity,
			&product.Curated,
			&product.ImageUrl,
		)
		if err != nil {
			log.Println(err)
		}
		products = append(products, product)
	}
	bytes, err := json.Marshal(products)
	internal.RenderJSON(w, bytes, http.StatusOK)
}
