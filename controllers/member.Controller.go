package controllers

import (
	"encoding/json"
	"membership-fitness-centre/models"
	"membership-fitness-centre/services"
	"net/http"
	"strconv"
)

type MemberController struct {
	service *services.MemberService
}

func NewMemberController(service *services.MemberService) *MemberController {
	return &MemberController{service: service}
}

func (uc *MemberController) CreateMember(w http.ResponseWriter, r *http.Request) {
	var member models.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uc.service.CreateMember(member.Name, member.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (uc *MemberController) GetMembers(w http.ResponseWriter, r *http.Request) {
	members, err := uc.service.GetMembers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(members)
}

func (uc *MemberController) UpdateMember(w http.ResponseWriter, r *http.Request) {
	var member models.Member

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = uc.service.UpdateMember(id, member.Name, member.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (uc *MemberController) DeleteMember(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = uc.service.DeleteMember(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
