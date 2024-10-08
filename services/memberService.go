package services

import (
	"database/sql"
	"errors"
	"membership-fitness-centre/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var mySigningKey = []byte("secret")

type MemberService struct {
	db *sql.DB
}

func NewMemberService(db *sql.DB) *MemberService {
	return &MemberService{db: db}
}

func (s *MemberService) CreateMember(username, email, password string) (string, string, error) {
	var exists bool
	ID := "0"

	query := `
        SELECT EXISTS (
            SELECT 1 FROM members WHERE username = $1 OR email = $2
        );`
	err := s.db.QueryRow(query, username, email).Scan(&exists)
	if err != nil {
		return "", ID, err
	}
	if exists {
		return "", ID, errors.New("member already exists")
	}

	query = `INSERT INTO members (username, email, password) VALUES ($1, $2, $3) RETURNING id`

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ID, err
	}

	err = s.db.QueryRow(query, username, email, string(bytes)).Scan(&ID)
	if err != nil {
		return "", ID, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":  ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", ID, err
	}

	return tokenString, ID, nil
}

func (s *MemberService) Authenticate(identifier, password string) (string, error) {
	var member models.Member

	query := `SELECT id, username, email, password FROM members WHERE username = $1 OR email = $1`
	err := s.db.QueryRow(query, identifier).Scan(&member.ID, &member.Username, &member.Email, &member.Password)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(password))
	if err != nil {
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
	bytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `UPDATE members SET password = $1 WHERE id = $2`
	_, err = s.db.Exec(query, string(bytes), memberID)
	return err
}
