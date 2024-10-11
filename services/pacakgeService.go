package services

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type PackageService struct {
	db *sql.DB
}

func NewPackageService(db *sql.DB) *PackageService {
	return &PackageService{db: db}
}

func (s *PackageService) AddPackage(memberID, packageID int) error {
	var interval int
	query := `select interval from packages where id = $1`
	err := s.db.QueryRow(query, packageID).Scan(&interval)
	if err != nil {
		return errors.New("package don't exists")
	}

	intervalQuery := strconv.Itoa(interval) + " month"
	query = fmt.Sprintf(`UPDATE members SET expiredate = NOW() + INTERVAL '%s', packageid = $1 WHERE id = $2`, intervalQuery)
	_, err = s.db.Exec(query, packageID, memberID)
	return err
}
