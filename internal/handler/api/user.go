package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/SemmiDev/kanbanapp/internal/entity"
	"github.com/SemmiDev/kanbanapp/internal/service"
)

type UserAPI interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	Delete(w http.ResponseWriter, r *http.Request)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.UserLogin

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("email or password is empty"))
		return
	}

	userArg := entity.User{
		Email:    user.Email,
		Password: user.Password,
	}

	userID, err := u.userService.Login(r.Context(), &userArg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	userResponse := map[string]interface{}{
		"user_id": userID,
		"message": "login success",
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  fmt.Sprintf("%d", userID),
		Path:   "/",
		MaxAge: 31536000,
		Domain: "",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResponse)
}

func (u *userAPI) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Fullname == "" || user.Email == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("register data is empty"))
		return
	}

	userArg := entity.User{
		Fullname: user.Fullname,
		Email:    user.Email,
		Password: user.Password,
	}

	userResult, err := u.userService.Register(r.Context(), &userArg)
	if err != nil {
		if err.Error() == "email already exists" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.NewErrorResponse(err.Error()))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	userResponse := map[string]interface{}{
		"user_id": float64(userResult.ID), // because json encode int to float
		"message": "register success",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse)
}

func (u *userAPI) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Value:  "",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "logout success"})
}

func (u *userAPI) Delete(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("user_id is empty"))
		return
	}

	deleteUserId, _ := strconv.Atoi(userId)

	err := u.userService.Delete(r.Context(), int(deleteUserId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "delete success"})
}
