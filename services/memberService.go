package services

import (
	"database/sql"
	"membership-fitness-centre/models"
)

type MemberService struct {
	db *sql.DB
}

func NewMemberService(db *sql.DB) *MemberService {
	return &MemberService{db: db}
}

func (s *MemberService) CreateMember(name, email string) (int, error) {
	sqlStatement := `INSERT INTO members (name, email) VALUES ($1, $2) RETURNING id`
	id := 0
	err := s.db.QueryRow(sqlStatement, name, email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *MemberService) GetMembers() ([]models.Member, error) {
	rows, err := s.db.Query("SELECT id, name, email FROM members")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.Member
	for rows.Next() {
		var member models.Member
		err = rows.Scan(&member.ID, &member.Name, &member.Email)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (s *MemberService) UpdateMember(id int, name, email string) error {
	sqlStatement := `UPDATE members SET name = $2, email = $3 WHERE id = $1`
	_, err := s.db.Exec(sqlStatement, id, name, email)
	return err
}

func (s *MemberService) DeleteMember(id int) error {
	sqlStatement := `DELETE FROM members WHERE id = $1`
	_, err := s.db.Exec(sqlStatement, id)
	return err
}
