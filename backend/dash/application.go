package dash

import (
	"time"
)

type Identifier struct {
	ID string `json:"id" db:"id"`
}

type ApplicationInfo struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	URL         string `json:"url" db:"url"`
	IconURL     string `json:"iconURL" db:"icon_url"`
}

// ApplicationRequest represents information about an application that would be submitted by a
// client of the API. (e.g. to create a new application).
type ApplicationRequest struct {
	ApplicationInfo
}

// ToApplicationRecord converts the ApplicationRequest to an ApplicationRecord. The returned
// ApplicationRecord will use the given ID and include updated timestamps.
func (a *ApplicationRequest) ToApplicationRecord(id string) *ApplicationRecord {
	now := time.Now().UTC()
	return &ApplicationRecord{
		Identifier: Identifier{ID: id},
		ApplicationInfo: ApplicationInfo{
			Name:        a.Name,
			Description: a.Description,
			URL:         a.URL,
			IconURL:     a.IconURL,
		},
		Created:  now,
		Modified: now,
	}
}

// ApplicationResponse represents information about an application that will be sent as a response
// to clients of the API.
type ApplicationResponse struct {
	Identifier
	ApplicationInfo
}

// ApplicationRecord represents an application that is stored in the database.
type ApplicationRecord struct {
	Identifier
	ApplicationInfo
	Created  time.Time `db:"created"`
	Modified time.Time `db:"modified"`
}

// ToApplicationResponse converts the ApplicationRecord to an ApplicationResponse to be returned by
// the API.
func (a *ApplicationRecord) ToApplicationResponse() *ApplicationResponse {
	return &ApplicationResponse{
		Identifier:      a.Identifier,
		ApplicationInfo: a.ApplicationInfo,
	}
}
