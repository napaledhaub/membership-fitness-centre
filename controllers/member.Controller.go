package controllers

import (
	"encoding/json"
	"membership-fitness-centre/models"
	"membership-fitness-centre/services"
	"net/http"
)

type MemberController struct {
	serviceMember *services.MemberService
}

func NewMemberController(serviceMember *services.MemberService) *MemberController {
	return &MemberController{serviceMember: serviceMember}
}

func (c *MemberController) CreateMember(w http.ResponseWriter, r *http.Request) {
	var member models.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := c.serviceMember.CreateMember(member.Username, member.Email, member.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (c *MemberController) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	token, err := c.serviceMember.Authenticate(creds.Identifier, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (c *MemberController) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID          int    `json:"ID"`
		NewPassword string `json:"new_password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	err := c.serviceMember.UpdatePassword(req.ID, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
