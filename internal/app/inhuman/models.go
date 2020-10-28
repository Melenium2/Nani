package inhuman

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type App struct {
	Id                int64             `json:"-" db:"id"`
	Bundle            string            `json:"bundle" db:"bundle"`
	DeveloperId       string            `json:"developerId" db:"developer_id"`
	Developer         string            `json:"developer" db:"developer"`
	Title             string            `json:"title" db:"title"`
	Categories        string            `json:"categories" db:"categories"`
	Price             string            `json:"price" db:"price"`
	Picture           string            `json:"picture" db:"picture"`
	Screenshots       []string          `json:"screenshots" db:"screenshots"`
	Rating            string            `json:"rating" db:"rating"`
	ReviewCount       string            `json:"reviewCount" db:"review_count"`
	RatingHistogram   []string          `json:"ratingHistogram" db:"rating_histogram"`
	Description       string            `json:"description" db:"description"`
	ShortDescription  string            `json:"shortDescription" db:"short_description"`
	RecentChanges     string            `json:"recentChanges" db:"recent_changes"`
	ReleaseDate       string            `json:"releaseDate" db:"release_date"`
	LastUpdateDate    string            `json:"lastUpdateDate" db:"last_update_date"`
	AppSize           string            `json:"appSize" db:"app_size"`
	Installs          string            `json:"installs" db:"installs"`
	Version           string            `json:"version" db:"version"`
	AndroidVersion    string            `json:"androidVersion" db:"android_version"`
	ContentRating     string            `json:"contentRating" db:"content_rating"`
	DeveloperContacts DeveloperContacts `json:"developerContacts" db:"developer_contacts"`
	PrivacyPolicy     string            `json:"privacyPolicy,omitempty"`
}

func (a App) Fields() string {
	return StructFields(reflect.TypeOf(a))[4:]
}

func (a App) String() string {
	screens, _ := json.Marshal(map[string]interface{}{
		"screenshots": a.Screenshots,
	})
	histogram, _ := json.Marshal(map[string]interface{}{
		"ratingHistogram": a.RatingHistogram,
	})
	contacts, _ := json.Marshal(map[string]interface{}{
		"developerContacts": a.DeveloperContacts,
	})

	return fmt.Sprintf(
		"%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s",
		a.Bundle,
		a.DeveloperId,
		a.Developer,
		a.Title,
		a.Categories,
		a.Price,
		a.Picture,
		string(screens),
		a.Rating,
		a.ReviewCount,
		string(histogram),
		a.Description,
		a.ShortDescription,
		a.RecentChanges,
		a.ReleaseDate,
		a.LastUpdateDate,
		a.AppSize,
		a.Installs,
		a.Version,
		a.AndroidVersion,
		a.ContentRating,
		string(contacts),
		a.PrivacyPolicy,
	)
}

type DeveloperContacts struct {
	Email    string `json:"email,omitempty"`
	Contacts string `json:"contacts,omitempty"`
}

type Keywords map[string]int
