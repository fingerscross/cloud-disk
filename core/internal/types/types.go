// Code generated by goctl. DO NOT EDIT.
package types

type LoginResponse struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Token string `json:"token"`
}
