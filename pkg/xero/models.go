package xero

type User struct {
	Id        string `json:"UserID"`
	Email     string `json:"EmailAddress"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Role      string `json:"OrganisationRole"`
}

type Organization struct {
	Id      string `json:"OrganisationID"`
	Name    string `json:"Name"`
	Country string `json:"CountryCode"`
}
