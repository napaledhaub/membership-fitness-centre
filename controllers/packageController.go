package controllers

import (
	"encoding/json"
	"membership-fitness-centre/services"
	"net/http"
)

type PackageController struct {
	servicePackage *services.PackageService
}

func NewPackageController(servicePackage *services.PackageService) *PackageController {
	return &PackageController{servicePackage: servicePackage}
}

func (c *PackageController) AddPackage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MemberID  int `json:"MemberID"`
		PackageID int `json:"PackageID"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	err := c.servicePackage.AddPackage(req.MemberID, req.PackageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
