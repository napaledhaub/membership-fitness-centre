package services

import (
	"database/sql"
	"errors"
	"membership-fitness-centre/models"
	"membership-fitness-centre/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var mySigningKey = []byte("secret")

type MemberService struct {
	db *sql.DB
}

func NewMemberService(db *sql.DB) *MemberService {
	return &MemberService{db: db}
}

func (s *MemberService) CreateMember(username, email, password string) (string, error) {
	var exists bool
	ID := "0"

	query := `
        SELECT EXISTS (
            SELECT 1 FROM members WHERE username = $1 OR email = $2
        );`
	err := s.db.QueryRow(query, username, email).Scan(&exists)
	if err != nil {
		return ID, err
	}
	if exists {
		return ID, errors.New("member already exists")
	}

	token := uuid.NewString()
	link := "http://localhost:8080/verify?token=" + token
	msg := []byte("To: " + email + "\n" +
		"Subject: Email Verification\n\n" +
		"Click the link to verify your email:" + link + token)
	err = utils.SendPackageEmails(email, msg)
	if err != nil {
		return ID, err
	}

	query = `INSERT INTO members (username, email, password, isverified, verificationtoken, tokencreatedat)
		VALUES ($1, $2, $3, false, $4, now()) RETURNING id`

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ID, err
	}

	err = s.db.QueryRow(query, username, email, string(bytes), token).Scan(&ID)
	if err != nil {
		return ID, err
	}

	return ID, nil
}

func (s *MemberService) Authenticate(identifier, password string) (string, error) {
	var member models.Member

	query := `SELECT id, username, email, password, isverified
		FROM members
		WHERE username = $1 OR email = $1`
	err := s.db.QueryRow(query, identifier).Scan(&member.ID,
		&member.Username,
		&member.Email,
		&member.Password,
		&member.IsVerified,
	)
	if err != nil {
		return "", err
	}
	if !member.IsVerified {
		return "", errors.New("not verified")
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

func (s *MemberService) VerifyEmail(token string) (string, error) {
	var member models.Member
	query := `SELECT id, email, tokencreatedat
		FROM members
		WHERE verificationtoken = $1 AND isverified = false`

	err := s.db.QueryRow(query, token).Scan(&member.ID,
		&member.Email,
		&member.TokenCreatedAt,
	)
	if err != nil {
		return "", err
	}

	if time.Since(member.TokenCreatedAt) > 24*time.Hour {
		return "", errors.New("link expired")
	}

	query = `UPDATE members SET isverified = true, verificationtoken = '' WHERE ID = $1`
	_, err = s.db.Exec(query, member.ID)
	if err != nil {
		return "", err
	}

	loginToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":  member.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := loginToken.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
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
