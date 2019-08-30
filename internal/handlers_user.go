package internal

import (
	"encoding/json"
	"fmt"
	"github.com/fernandochristyanto/devcamp-backend/model"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

// GetUserByID a method to get user given userID params in URL
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	nowTime := time.Now()

	userID, err := strconv.ParseInt(param.ByName("userID"), 10, 64)
	if err != nil {
		log.Printf("[internal][GetUserByID] fail parsing userid  %s :%+v\n",
			param.ByName("userID"), err)
		renderJSON(w, []byte(
			`message: "Fail parsing user id"`,
		), http.StatusBadRequest)
		return
	}
	query := fmt.Sprintf("SELECT id, coalesce(name, '') FROM users WHERE id = %d ", userID)
	rows, err := h.DB.Query(query)
	if err != nil {
		log.Printf("[internal][GetUserByID] fail select user where user_id = %s :%+v\n",
			param.ByName("userID"), err)
		return
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		user := model.User{}
		err := rows.Scan(
			&user.ID,
			// &user.Name,
		)
		if err != nil {
			log.Printf("[internal][GetUserByID] fail scanning row %+v\n", err)
			continue
		}
		users = append(users, user)
	}
	defer rows.Close()

	bytes, err := json.Marshal(users)
	if err != nil {
		log.Printf("[internal][GetUserByID] fail marshalling json %+v\n", err)
		return
	}
	renderJSON(w, bytes, http.StatusOK)

	processTime := time.Since(nowTime).Seconds()
	log.Printf("Process time : %f seconds\n", processTime)

}

func (h *Handler) GetUserByEmailAndPassword(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	renderJSON(w, []byte(`
	{
		"id": 1,
		"email": "admin@admin"
	}`), http.StatusOK)
}
