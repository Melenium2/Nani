package db

import (
	"Nani/internal/app/config"
	"Nani/internal/app/inhuman"
	"context"
	"database/sql"
	"fmt"
	"github.com/mailru/go-clickhouse"
)

type AppRepository interface {
	Insert(ctx context.Context, app *inhuman.App) error
	InsertBatch(ctx context.Context, apps[] *inhuman.App) error
}

type ClickhouseDatabase struct {
	connection *sql.DB
}

func (c *ClickhouseDatabase) Insert(ctx context.Context, app *inhuman.App) error {
	_, err := c.connection.ExecContext(
		ctx,
		fmt.Sprintf("insert into apps (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", app.Fields()),
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
		clickhouse.Array([]string{ app.DeveloperContacts.Email }),
		clickhouse.Array([]string{ app.DeveloperContacts.Contacts }),
		app.PrivacyPolicy,
	)

	if err != nil {
		return err
	}

	return nil
}

func (c *ClickhouseDatabase) InsertBatch(ctx context.Context, apps []*inhuman.App) error {
	if len(apps) == 0 {
		return nil
	}

	t, err := c.connection.Begin()
	if err != nil {
		return err
	}
	stmt, err := t.PrepareContext(
		ctx,
		fmt.Sprintf("insert into apps (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", apps[0].Fields()),
	)
	if err != nil {
		return err
	}
	for _, v := range apps {
		_, err := stmt.ExecContext(ctx,
			v.Bundle,
			v.DeveloperId,
			v.Developer,
			v.Title,
			v.Categories,
			v.Price,
			v.Picture,
			clickhouse.Array(v.Screenshots),
			v.Rating,
			v.ReviewCount,
			clickhouse.Array(v.RatingHistogram),
			v.Description,
			v.ShortDescription,
			v.RecentChanges,
			v.ReleaseDate,
			v.LastUpdateDate,
			v.AppSize,
			v.Installs,
			v.Version,
			v.AndroidVersion,
			v.ContentRating,
			clickhouse.Array([]string{ v.DeveloperContacts.Email }),
			clickhouse.Array([]string{ v.DeveloperContacts.Contacts }),
			v.PrivacyPolicy,
		)
		if err != nil {
			return err
		}
	}

	if err := t.Commit(); err != nil {
		return err
	}

	return nil
}

func New(config config.DBConfig) *ClickhouseDatabase {
	if config.Connection == nil {
		url, err := ConnectionUrl(config)
		if err != nil {
			panic(err)
		}
		config.Connection, err = Connect(url)
		if err != nil {
			panic(err)
		}
	}

	if config.Schema != "" {
		err := InitSchema(config.Connection, config.Schema)
		if err != nil {
			panic(err)
		}
	}

	c := &ClickhouseDatabase{
		connection: config.Connection,
	}

	return c
}
