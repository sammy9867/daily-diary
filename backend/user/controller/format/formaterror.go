package format

import (
	"errors"
	"strings"

	"github.com/badoux/checkmail"
	"github.com/sammy9867/daily-diary/backend/domain"
)

// CheckError is used to check the type of error
func CheckError(err string) error {
	if strings.Contains(err, "username") {
		return errors.New("Username is already taken")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email is already taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}

	return errors.New("Incorrect Details")
}

//Validate is used to check if the user has entered correct input format
func Validate(u *domain.User, action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Username == "" {
			return errors.New("Username cannot be blank")
		}
		if u.Password == "" {
			return errors.New("Password cannot be blank")
		}
		if u.Email == "" {
			return errors.New("Email cannot be blank")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	case "login":
		if u.Password == "" {
			return errors.New("Password cannot be blank")
		}
		if u.Email == "" {
			return errors.New("Email cannot be blank")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Username == "" {
			return errors.New("Username cannot be blank")
		}
		if u.Password == "" {
			return errors.New("Password cannot be blank")
		}
		if u.Email == "" {
			return errors.New("Email cannot be blank")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}
