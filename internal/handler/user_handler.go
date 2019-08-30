package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/fernandochristyanto/devcamp-backend/internal"
	"github.com/fernandochristyanto/devcamp-backend/model/dto"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
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
	if sellerRegistrationDTO.UserId == nil {
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

	internal.RenderJSON(w, []byte(`
	{
		"message": "Success"
	}`), http.StatusOK)
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
