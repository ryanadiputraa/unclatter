package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	uuid := uuid.NewString()

	cases := []struct {
		name string
		arg  struct {
			email     string
			firstName string
			lastName  string
		}
		expected *User
		error
	}{
		{
			name: "should return a valid user",
			arg: struct {
				email     string
				firstName string
				lastName  string
			}{
				email:     "test@mail.com",
				firstName: "Test",
				lastName:  "Lastname",
			},
			expected: &User{
				ID:        uuid,
				Email:     "test@mail.com",
				FirstName: "Test",
				LastName:  "Lastname",
				CreatedAt: time.Now().UTC(),
			},
			error: nil,
		},
		{
			name: "should return error when given invalid email address",
			arg: struct {
				email     string
				firstName string
				lastName  string
			}{
				email:     "invalidemailaddress",
				firstName: "Test",
				lastName:  "Lastname",
			},
			expected: nil,
			error:    errors.New("invalid email address"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			user, err := NewUser(c.arg.email, c.arg.firstName, c.arg.lastName)
			assert.Equal(t, c.error, err)
			if err != nil {
				return
			}

			assert.NotEmpty(t, user.ID)
			assert.Equal(t, c.expected.Email, user.Email)
			assert.Equal(t, c.expected.FirstName, user.FirstName)
			assert.Equal(t, c.expected.LastName, user.LastName)
			assert.NotEmpty(t, user.CreatedAt)
		})
	}
}
