package responses

import (
	"time"
)

type Platform string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
)

type IOSCred struct {
	Platform            Platform               `json:"Platform"`
	Label               string                 `json:"label"`
	Id                  string                 `json:"credentialid"`
	Created             time.Time              `json:"created"`
	LastMod             time.Time              `json:"lastMod"`
	Certificate         IOSCert                `json:"certificate"`
	ProvisioningProfile IOSProvisioningProfile `json:"provisioningProfile"`
	Links               map[string]Link        `json:"links"`
}

type IOSCert struct {
	TeamId         string    `json:"teamId"`
	Name           string    `json:"certName"`
	Expiration     time.Time `json:"expiration"`
	IsDistribution bool      `json:"isDistribution"`
	Uploaded       string    `json:"uploaded"`
}

type IOSProvisioningProfile struct {
	TeamID              string    `json:"teamId"`
	BundleID            string    `json:"bundleId"`
	Expiration          time.Time `json:"expiration"`
	IsEnterpriseProfile bool      `json:"isEnterpriseProfile"`
	Type                string    `json:"type"`
	NumDevices          int       `json:"numDevices"`
}

type Link struct {
	Method string `json:"method"`
	Href   string `json:"href"`
}
