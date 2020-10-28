package db_test

import (
	"Nani/internal/app/db"
	"Nani/internal/app/inhuman"
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mailru/go-clickhouse"
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

func Db() (*sql.DB, func()) {
	config := db.Config{
		Database: "default",
		Address:  "192.168.99.100",
		Port:     "8123",
		Schema:   "",
	}
	url, _ := db.ConnectionUrl(config)
	driver, _ := db.Connect(url)
	db.InitSchema(driver, "../../../config/schema.sql")
	return driver, func() {
		driver.Exec("drop table apps")
	}
}

func TestInsertMock_ShouldInsertNewRecordToDb_NoError(t *testing.T) {
	ctx := context.Background()
	d, mock := MockDb()
	defer d.Close()

	app := App()
	str := fmt.Sprintf("^insert into apps \\(.+\\) values \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?\\)$")
	mock.ExpectExec(str).
		WithArgs(
			app.Bundle,
			app.DeveloperId,
			app.Developer,
			app.Title,
			app.Categories,
			app.Price,
			app.Picture,
			clickhouse.Array(app.Screenshots),
			app.Rating,
			app.ReviewCount,
			clickhouse.Array(app.RatingHistogram),
			app.Description,
			app.ShortDescription,
			app.RecentChanges,
			app.ReleaseDate,
			app.LastUpdateDate,
			app.AppSize,
			app.Installs,
			app.Version,
			app.AndroidVersion,
			app.ContentRating,
			clickhouse.Array([]string{app.DeveloperContacts.Email}),
			clickhouse.Array([]string{app.DeveloperContacts.Contacts}),
			app.PrivacyPolicy).
		WillReturnResult(sqlmock.NewErrorResult(nil))

	var rep db.Repository
	assert.NotPanics(t, func() {
		rep = db.New(db.Config{Connection: d})
	})
	assert.NoError(t, rep.Insert(ctx, app))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertBatchMock_ShouldInsertNewRecords_NoError(t *testing.T) {
	ctx := context.Background()
	d, mock := MockDb()
	defer d.Close()

	apps := []*inhuman.App { App(), App(), App() }
	mock.ExpectBegin()
	str := fmt.Sprintf("^insert into apps \\(.+\\) values \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?\\)$")
	stmt := mock.ExpectPrepare(str)
	for _, app := range apps {
		stmt.ExpectExec().
			WithArgs(
				app.Bundle,
				app.DeveloperId,
				app.Developer,
				app.Title,
				app.Categories,
				app.Price,
				app.Picture,
				clickhouse.Array(app.Screenshots),
				app.Rating,
				app.ReviewCount,
				clickhouse.Array(app.RatingHistogram),
				app.Description,
				app.ShortDescription,
				app.RecentChanges,
				app.ReleaseDate,
				app.LastUpdateDate,
				app.AppSize,
				app.Installs,
				app.Version,
				app.AndroidVersion,
				app.ContentRating,
				clickhouse.Array([]string{app.DeveloperContacts.Email}),
				clickhouse.Array([]string{app.DeveloperContacts.Contacts}),
				app.PrivacyPolicy).
			WillReturnResult(sqlmock.NewErrorResult(nil))
	}
	mock.ExpectCommit()

	var rep db.Repository
	assert.NotPanics(t, func() {
		rep = db.New(db.Config{Connection: d})
	})
	assert.NoError(t, rep.InsertBatch(ctx, apps))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertBatchMock_ShouldRiseErrorCozDbConnectionLost_Error(t *testing.T) {
	ctx := context.Background()
	d, mock := MockDb()

	apps := []*inhuman.App { App(), App(), App() }
	mock.ExpectBegin()
	str := fmt.Sprintf("^insert into apps \\(.+\\) values \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?, \\?\\)$")
	stmt := mock.ExpectPrepare(str)
	for _, app := range apps {
		stmt.ExpectExec().
			WithArgs(
				app.Bundle,
				app.DeveloperId,
				app.Developer,
				app.Title,
				app.Categories,
				app.Price,
				app.Picture,
				clickhouse.Array(app.Screenshots),
				app.Rating,
				app.ReviewCount,
				clickhouse.Array(app.RatingHistogram),
				app.Description,
				app.ShortDescription,
				app.RecentChanges,
				app.ReleaseDate,
				app.LastUpdateDate,
				app.AppSize,
				app.Installs,
				app.Version,
				app.AndroidVersion,
				app.ContentRating,
				clickhouse.Array([]string{app.DeveloperContacts.Email}),
				clickhouse.Array([]string{app.DeveloperContacts.Contacts}),
				app.PrivacyPolicy).
			WillReturnResult(sqlmock.NewErrorResult(nil))
	}
	d.Close()
	mock.ExpectCommit()

	var rep db.Repository
	assert.NotPanics(t, func() {
		rep = db.New(db.Config{Connection: d})
	})
	assert.Error(t, rep.InsertBatch(ctx, apps))

	assert.Error(t, mock.ExpectationsWereMet())
}


func TestInsert_ShouldInsertNewRow_NoError(t *testing.T) {
	ctx := context.Background()
	driver, cleaner := Db()
	defer cleaner()

	repo := db.New(db.Config{ Connection: driver})
	err := repo.Insert(ctx, App())
	assert.NoError(t, err)

	var bundle string
	err = driver.QueryRowContext(ctx, fmt.Sprint("select bundle from apps")).
		Scan(&bundle)
	assert.NoError(t, err)
	assert.NotEmpty(t, bundle)
}

func TestInsert_ShouldReturnErrorCozDriverClosed_Error(t *testing.T) {
	ctx := context.Background()
	driver, cleaner := Db()
	defer cleaner()

	repo := db.New(db.Config{ Connection: driver})
	driver.Close()
	err := repo.Insert(ctx, App())
	assert.Error(t, err)

	var bundle string
	err = driver.QueryRowContext(ctx, fmt.Sprint("select bundle from apps")).
		Scan(&bundle)
	assert.Error(t, err)
	assert.Empty(t, bundle)
}

func TestInsertBatch_ShouldInsertSomeItems_NoError(t *testing.T) {
	ctx := context.Background()
	driver, cleaner := Db()
	defer cleaner()

	apps := []*inhuman.App{ App(), App(), App() }
	repo := db.New(db.Config{ Connection: driver})
	err := repo.InsertBatch(ctx, apps)
	assert.NoError(t, err)

	var bundle string
	rows, err := driver.QueryContext(ctx, fmt.Sprint("select bundle from apps"))
	assert.NoError(t, err)
	assert.NotNil(t, rows)

	var counter = 0
	for rows.Next() {
		assert.NoError(t, rows.Scan(&bundle))
		counter++
	}

	assert.Equal(t, 3, counter)
}


