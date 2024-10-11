package services

import (
	"database/sql"
	"membership-fitness-centre/models"
	"time"
)

type Subscriptionervice struct {
	DB *sql.DB
}

func NewSubscriptionervice(db *sql.DB) *Subscriptionervice {
	return &Subscriptionervice{DB: db}
}

func (s *Subscriptionervice) GetExpiringSubscriptions() ([]models.Member, error) {
	today := time.Now()
	tomorrow := today.Add(48 * time.Hour)

	query := `
        SELECT id, email, expiredate
        FROM members
        WHERE expiredate BETWEEN $1 AND $2`

	rows, err := s.DB.Query(query, today, tomorrow)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.Member
	for rows.Next() {
		var mem models.Member
		if err := rows.Scan(&mem.ID, &mem.Email, &mem.ExpireDate); err != nil {
			return nil, err
		}
		members = append(members, mem)
	}
	return members, nil
}
