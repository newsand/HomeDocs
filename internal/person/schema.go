package person

import (
	"errors"
	"strings"
)

var ErrNameRequired = errors.New("name required")
var ErrNotFound = errors.New("person not found")

func ValidName(p_Name string) bool {
	return strings.TrimSpace(p_Name) != ""
}

type Upload struct {
	Name  string
	Bytes []byte
}

type CNH struct {
	RegisterNumber string
	Category       string
	ExpiresOn      string
	FirstLicenseOn string
	IssuedOn       string
	Notes          string
}

type CPF struct {
	Number       string
	BirthDate    string
	Status       string
	RegisteredOn string
}

type RG struct {
	Number      string
	SocialName  string
	BirthPlace  string
	FatherName  string
	MotherName  string
	IssuedOn    string
	Issuer      string
	IssuerState string
	ExpiresOn   string
	Sex         string
}

type VoterID struct {
	Number   string
	Zone     string
	Section  string
	City     string
	State    string
	IssuedOn string
}

type Passport struct {
	Kind           string
	IssuingCountry string
	Number         string
	Nationality    string
	IssuedOn       string
	Authority      string
	ExpiresOn      string
}

type Person struct {
	ID            string
	Name          string
	ProfilePhoto  string
	DocumentScans []string
	CNH           CNH
	CPF           CPF
	RG            RG
	VoterID       VoterID
	Passport      Passport
}

type Input struct {
	Name          string
	ProfilePhoto  *Upload
	DocumentScans []Upload
	CNH           CNH
	CPF           CPF
	RG            RG
	VoterID       VoterID
	Passport      Passport
}
