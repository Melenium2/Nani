package inhuman_test

import (
	"Nani/internal/app/inhuman"
	"testing"
)

func TestAppString_ShouldReturnStringRepresentationOfFields_NoError(t *testing.T) {
	app := &inhuman.App{
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

	str := app.Fields()
	t.Log(str)
}
