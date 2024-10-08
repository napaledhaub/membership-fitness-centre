package services

import (
	"database/sql"
	"errors"
	"membership-fitness-centre/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("secret")

type MemberService struct {
	db *sql.DB
}

func NewMemberService(db *sql.DB) *MemberService {
	return &MemberService{db: db}
}

func (s *MemberService) CreateMember(username, email, password string) (string, error) {
	sqlStatement := `INSERT INTO members (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	id := 0
	err := s.db.QueryRow(sqlStatement, username, email, password).Scan(&id)
	if err != nil {
		return "", err
	}

	token, err := s.Authenticate(username, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *MemberService) Authenticate(identifier, password string) (string, error) {
	var member models.Member

	query := `SELECT id, username, email, password FROM members WHERE username = $1 OR email = $1`
	err := s.db.QueryRow(query, identifier).Scan(&member.ID, &member.Username, &member.Email, &member.Password)
	if err != nil {
		return "", err
	}

	if password != member.Password {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":  member.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *MemberService) UpdatePassword(memberID int, newPassword string) error {
	query := `UPDATE members SET password = $1 WHERE id = $2`
	_, err := s.db.Exec(query, newPassword, memberID)
	return err
}
