package person

func columnNames() []string {
	return []string{
		"id",
		"name",
		"profile_photo",
		"cnh_register_number",
		"cnh_category",
		"cnh_expires_on",
		"cnh_first_license_on",
		"cnh_issued_on",
		"cnh_notes",
		"cpf_number",
		"cpf_birth_date",
		"cpf_status",
		"cpf_registered_on",
		"rg_number",
		"rg_social_name",
		"rg_birth_place",
		"rg_father_name",
		"rg_mother_name",
		"rg_issued_on",
		"rg_issuer",
		"rg_issuer_state",
		"rg_expires_on",
		"rg_sex",
		"voter_number",
		"voter_zone",
		"voter_section",
		"voter_city",
		"voter_state",
		"voter_issued_on",
		"passport_kind",
		"passport_issuing_country",
		"passport_number",
		"passport_nationality",
		"passport_issued_on",
		"passport_authority",
		"passport_expires_on",
	}
}

func searchColumns() []string {
	names := columnNames()
	out := make([]string, 0, len(names)-2)
	for _, name := range names {
		if name == "id" || name == "profile_photo" {
			continue
		}
		out = append(out, name)
	}
	return out
}

func columnValues(p_Person Person) []any {
	return []any{
		p_Person.ID,
		p_Person.Name,
		p_Person.ProfilePhoto,
		p_Person.CNH.RegisterNumber,
		p_Person.CNH.Category,
		p_Person.CNH.ExpiresOn,
		p_Person.CNH.FirstLicenseOn,
		p_Person.CNH.IssuedOn,
		p_Person.CNH.Notes,
		p_Person.CPF.Number,
		p_Person.CPF.BirthDate,
		p_Person.CPF.Status,
		p_Person.CPF.RegisteredOn,
		p_Person.RG.Number,
		p_Person.RG.SocialName,
		p_Person.RG.BirthPlace,
		p_Person.RG.FatherName,
		p_Person.RG.MotherName,
		p_Person.RG.IssuedOn,
		p_Person.RG.Issuer,
		p_Person.RG.IssuerState,
		p_Person.RG.ExpiresOn,
		p_Person.RG.Sex,
		p_Person.VoterID.Number,
		p_Person.VoterID.Zone,
		p_Person.VoterID.Section,
		p_Person.VoterID.City,
		p_Person.VoterID.State,
		p_Person.VoterID.IssuedOn,
		p_Person.Passport.Kind,
		p_Person.Passport.IssuingCountry,
		p_Person.Passport.Number,
		p_Person.Passport.Nationality,
		p_Person.Passport.IssuedOn,
		p_Person.Passport.Authority,
		p_Person.Passport.ExpiresOn,
	}
}

func scanDest(p_Person *Person) []any {
	return []any{
		&p_Person.ID,
		&p_Person.Name,
		&p_Person.ProfilePhoto,
		&p_Person.CNH.RegisterNumber,
		&p_Person.CNH.Category,
		&p_Person.CNH.ExpiresOn,
		&p_Person.CNH.FirstLicenseOn,
		&p_Person.CNH.IssuedOn,
		&p_Person.CNH.Notes,
		&p_Person.CPF.Number,
		&p_Person.CPF.BirthDate,
		&p_Person.CPF.Status,
		&p_Person.CPF.RegisteredOn,
		&p_Person.RG.Number,
		&p_Person.RG.SocialName,
		&p_Person.RG.BirthPlace,
		&p_Person.RG.FatherName,
		&p_Person.RG.MotherName,
		&p_Person.RG.IssuedOn,
		&p_Person.RG.Issuer,
		&p_Person.RG.IssuerState,
		&p_Person.RG.ExpiresOn,
		&p_Person.RG.Sex,
		&p_Person.VoterID.Number,
		&p_Person.VoterID.Zone,
		&p_Person.VoterID.Section,
		&p_Person.VoterID.City,
		&p_Person.VoterID.State,
		&p_Person.VoterID.IssuedOn,
		&p_Person.Passport.Kind,
		&p_Person.Passport.IssuingCountry,
		&p_Person.Passport.Number,
		&p_Person.Passport.Nationality,
		&p_Person.Passport.IssuedOn,
		&p_Person.Passport.Authority,
		&p_Person.Passport.ExpiresOn,
	}
}
