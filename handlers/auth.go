package handlers

import (
	"encoding/json"
	"fileupbackendv2/internal/auth"
	"fileupbackendv2/pkg/response"
	"net/http"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var signupRequest SignupRequest

	err := json.NewDecoder(r.Body).Decode(&signupRequest)
	if err != nil {
		response.BadRequestResponse(w, err.Error())
		return
	}

	err = auth.CreateUser(signupRequest.Username, signupRequest.Email, signupRequest.Password)
	if err != nil {
		response.InternalServerErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, "User created successfully")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var signupRequest SignupRequest

	err := json.NewDecoder(r.Body).Decode(&signupRequest)
	if err != nil {
		response.BadRequestResponse(w, err.Error())
		return
	}
	if signupRequest.Username == "" && signupRequest.Email == "" {
		response.BadRequestResponse(w, "Username or email is required")
		return
	}
	if signupRequest.Password == "" {
		response.BadRequestResponse(w, "Password is required")
		return
	}
	token, err := auth.LoginUser(signupRequest.Username, signupRequest.Email, signupRequest.Password)
	if err != nil {
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	response.JSONResponse(w, map[string]string{
		"token": token,
	})

}

func GenerateAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "API Key generated successfully",
	})
}
