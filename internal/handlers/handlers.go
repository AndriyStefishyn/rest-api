package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rest-api/internal/service"
	"rest-api/internal/shop"

	"github.com/gorilla/mux"
)

type ShopHandlers struct {
	ShopService *service.ShopService
}

func NewShopHandlers(ss *service.ShopService) *ShopHandlers {
	return &ShopHandlers{ShopService: ss}
}

func (sh *ShopHandlers) GetShopByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	shop, err := sh.ShopService.GetShop(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(shop)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)

}

func readBody(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return []byte{}, err
	}

	defer r.Body.Close()

	return rBody, nil
}
func (sh *ShopHandlers) GetShopsHandler(w http.ResponseWriter, r *http.Request) {
	shops, err := sh.ShopService.GetAllShops(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(shops)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (sh *ShopHandlers) CreateShopHandler(w http.ResponseWriter, r *http.Request) {
	rBody, err := readBody(w, r)
	if err != nil {
		return
	}

	var shop shop.Shop

	err = json.Unmarshal(rBody, &shop)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := sh.ShopService.CreateShop(r.Context(), shop)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("shops was created with id :%s", id)
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
}

func (sh *ShopHandlers) DeleteShopHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := sh.ShopService.DeleteShop(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("shop with id :%s was successfully deleted", id)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (sh *ShopHandlers) UpdateShopHandler(w http.ResponseWriter, r *http.Request) {
	rBody, err := readBody(w, r)
	if err != nil {
		return
	}

	var updates shop.Shop

	err = json.Unmarshal(rBody, &updates)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = sh.ShopService.UpdateShop(r.Context(), updates)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode("shop was successfully updated")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
