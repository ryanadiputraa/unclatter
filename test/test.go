package test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ryanadiputraa/unclatter/app/article"
	"github.com/ryanadiputraa/unclatter/app/auth"
	"github.com/ryanadiputraa/unclatter/app/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	TestUser = &user.User{
		ID:        uuid.NewString(),
		Email:     "johndoe@mail.com",
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now().UTC(),
	}
	TestAuthProvider = &auth.AuthProvider{
		ID:             uuid.NewString(),
		Provider:       "google",
		ProviderUserID: "129301293801231",
		UserID:         TestUser.ID,
		CreatedAt:      time.Now().UTC(),
	}
	TestArticle = &article.Article{
		ID:          uuid.NewString(),
		Title:       "Title",
		Content:     "<div><a onblur=\"alert(secret)\" href=\"http://www.google.com\">Google</a><p>article content</p></div>",
		ArticleLink: "https://unclatter.com",
		UserID:      TestUser.ID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	TestArticle2 = &article.Article{
		ID:          uuid.NewString(),
		Title:       "Title 2",
		Content:     "<div><a onblur=\"alert(secret)\" href=\"http://www.google.com\">Google</a><p>article content 2</p></div>",
		ArticleLink: "https://unclatter.com/2",
		UserID:      TestUser.ID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	TestArticle3 = &article.Article{
		ID:          uuid.NewString(),
		Title:       "Title 3",
		Content:     "<div><a onblur=\"alert(secret)\" href=\"http://www.google.com\">Google</a><p>article content 3</p></div>",
		ArticleLink: "https://unclatter.com/3",
		UserID:      TestUser.ID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
)

func NewMockDB(t *testing.T) (*gorm.DB, *sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("fail to create mock db conn: ", err.Error())
	}

	gormDB, err := gorm.Open(postgres.New(
		postgres.Config{
			Conn:       db,
			DriverName: "postgres",
		},
	), &gorm.Config{})
	if err != nil {
		t.Fatal("fail to open db conn: ", err.Error())
	}
	return gormDB, db, mock
}
