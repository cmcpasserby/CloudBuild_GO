package responses

import "time"

type Project struct {
	Name                 string          `json:"name"`
	Id                   string          `json:"projectId"`
	OrgName              string          `json:"OrgName"`
	Guid                 string          `json:"guid"`
	Created              time.Time       `json:"created"`
	Links                map[string]Link `json:"links"`
	Disabled             bool            `json:"disabled"`
	DisableNotifications bool            `json:"disableNotifications"`
	GenerateShareLinks   bool            `json:"generateShareLinks"`
}
