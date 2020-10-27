package db_test

import (
	"Nani/internal/app/db"
	"Nani/internal/app/inhuman"
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func App() *inhuman.App {
	return &inhuman.App{
		Bundle:            "com.ky",
		DeveloperId:       "devid",
		Developer:         "imdeveloper",
		Title:             "super title",
		Categories:        "GAME",
		Price:             "",
		Picture:           "http://localhost:8080",
		Screenshots:       []string{"1", "2", "3"},
		Rating:            "4.3",
		ReviewCount:       "10002",
		RatingHistogram:   []string{"515", "1323", "12333", "12323213", "0000000"},
		Description:       "hi",
		ShortDescription:  "super hi",
		RecentChanges:     "uuuyyyy",
		ReleaseDate:       "2020-06-20",
		LastUpdateDate:    "2020-07-29",
		AppSize:           "10M",
		Installs:          "100000+",
		Version:           "1.3.23",
		AndroidVersion:    "4.3+",
		ContentRating:     "12+",
		DeveloperContacts: inhuman.DeveloperContacts{Email: "email", Contacts: "conatasdsd"},
		PrivacyPolicy:     "http://localhost/hello",
	}
}

func MockDb() (*sql.DB, sqlmock.Sqlmock) {
	d, m, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	return d, m
}

func TestInsertMock_ShouldInsertNewRecordToDb_NoError(t *testing.T) {
	ctx := context.Background()
	d, mock := MockDb()
	defer d.Close()

	app := App()
	str := fmt.Sprintf("insert into apps \\(%s\\) values \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?\\)", app.Fields())
	mock.ExpectExec(str).
		WillReturnResult(sqlmock.NewErrorResult(nil))

	var rep db.Repository
	assert.NotPanics(t, func() {
		rep = db.New(db.Config{Connection: d})
	})
	assert.NoError(t, rep.Insert(ctx, app))

	assert.NoError(t, mock.ExpectationsWereMet())
}
