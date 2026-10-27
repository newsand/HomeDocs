package http

import "idcenter/internal/person"

type formField struct {
	label string
	name  string
	read  func(p_Person person.Person) string
	apply func(p_Input *person.Input, p_Value string)
}

type formGroup struct {
	title  string
	fields []formField
}

type ViewField struct {
	Label string
	Name  string
	Value string
}

type ViewGroup struct {
	Title  string
	Fields []ViewField
}

type ListItem struct {
	ID   string
	Name string
}

func catalog() []formGroup {
	return []formGroup{
		{title: "Dados", fields: []formField{
			text("Nome", "name", func(p person.Person) string { return p.Name }, func(in *person.Input, v string) { in.Name = v }),
		}},
		{title: "CNH", fields: []formField{
			text("Nº registro", "cnh_register_number", func(p person.Person) string { return p.CNH.RegisterNumber }, func(in *person.Input, v string) { in.CNH.RegisterNumber = v }),
			text("Categoria", "cnh_category", func(p person.Person) string { return p.CNH.Category }, func(in *person.Input, v string) { in.CNH.Category = v }),
			text("Validade", "cnh_expires_on", func(p person.Person) string { return p.CNH.ExpiresOn }, func(in *person.Input, v string) { in.CNH.ExpiresOn = v }),
			text("1ª habilitação", "cnh_first_license_on", func(p person.Person) string { return p.CNH.FirstLicenseOn }, func(in *person.Input, v string) { in.CNH.FirstLicenseOn = v }),
			text("Emissão", "cnh_issued_on", func(p person.Person) string { return p.CNH.IssuedOn }, func(in *person.Input, v string) { in.CNH.IssuedOn = v }),
			text("Observações", "cnh_notes", func(p person.Person) string { return p.CNH.Notes }, func(in *person.Input, v string) { in.CNH.Notes = v }),
		}},
		{title: "CPF", fields: []formField{
			text("Número", "cpf_number", func(p person.Person) string { return p.CPF.Number }, func(in *person.Input, v string) { in.CPF.Number = v }),
			text("Nascimento", "cpf_birth_date", func(p person.Person) string { return p.CPF.BirthDate }, func(in *person.Input, v string) { in.CPF.BirthDate = v }),
			text("Situação cadastral", "cpf_status", func(p person.Person) string { return p.CPF.Status }, func(in *person.Input, v string) { in.CPF.Status = v }),
			text("Data da inscrição", "cpf_registered_on", func(p person.Person) string { return p.CPF.RegisteredOn }, func(in *person.Input, v string) { in.CPF.RegisteredOn = v }),
		}},
		{title: "RG", fields: []formField{
			text("Número", "rg_number", func(p person.Person) string { return p.RG.Number }, func(in *person.Input, v string) { in.RG.Number = v }),
			text("Nome social", "rg_social_name", func(p person.Person) string { return p.RG.SocialName }, func(in *person.Input, v string) { in.RG.SocialName = v }),
			text("Naturalidade", "rg_birth_place", func(p person.Person) string { return p.RG.BirthPlace }, func(in *person.Input, v string) { in.RG.BirthPlace = v }),
			text("Pai", "rg_father_name", func(p person.Person) string { return p.RG.FatherName }, func(in *person.Input, v string) { in.RG.FatherName = v }),
			text("Mãe", "rg_mother_name", func(p person.Person) string { return p.RG.MotherName }, func(in *person.Input, v string) { in.RG.MotherName = v }),
			text("Expedição", "rg_issued_on", func(p person.Person) string { return p.RG.IssuedOn }, func(in *person.Input, v string) { in.RG.IssuedOn = v }),
			text("Órgão expedidor", "rg_issuer", func(p person.Person) string { return p.RG.Issuer }, func(in *person.Input, v string) { in.RG.Issuer = v }),
			text("UF", "rg_issuer_state", func(p person.Person) string { return p.RG.IssuerState }, func(in *person.Input, v string) { in.RG.IssuerState = v }),
			text("Validade", "rg_expires_on", func(p person.Person) string { return p.RG.ExpiresOn }, func(in *person.Input, v string) { in.RG.ExpiresOn = v }),
			text("Sexo", "rg_sex", func(p person.Person) string { return p.RG.Sex }, func(in *person.Input, v string) { in.RG.Sex = v }),
		}},
		{title: "Título de eleitor", fields: []formField{
			text("Inscrição", "voter_number", func(p person.Person) string { return p.VoterID.Number }, func(in *person.Input, v string) { in.VoterID.Number = v }),
			text("Zona", "voter_zone", func(p person.Person) string { return p.VoterID.Zone }, func(in *person.Input, v string) { in.VoterID.Zone = v }),
			text("Seção", "voter_section", func(p person.Person) string { return p.VoterID.Section }, func(in *person.Input, v string) { in.VoterID.Section = v }),
			text("Município", "voter_city", func(p person.Person) string { return p.VoterID.City }, func(in *person.Input, v string) { in.VoterID.City = v }),
			text("UF", "voter_state", func(p person.Person) string { return p.VoterID.State }, func(in *person.Input, v string) { in.VoterID.State = v }),
			text("Emissão", "voter_issued_on", func(p person.Person) string { return p.VoterID.IssuedOn }, func(in *person.Input, v string) { in.VoterID.IssuedOn = v }),
		}},
		{title: "Passaporte", fields: []formField{
			text("Tipo", "passport_kind", func(p person.Person) string { return p.Passport.Kind }, func(in *person.Input, v string) { in.Passport.Kind = v }),
			text("País emissor", "passport_issuing_country", func(p person.Person) string { return p.Passport.IssuingCountry }, func(in *person.Input, v string) { in.Passport.IssuingCountry = v }),
			text("Número", "passport_number", func(p person.Person) string { return p.Passport.Number }, func(in *person.Input, v string) { in.Passport.Number = v }),
			text("Nacionalidade", "passport_nationality", func(p person.Person) string { return p.Passport.Nationality }, func(in *person.Input, v string) { in.Passport.Nationality = v }),
			text("Emissão", "passport_issued_on", func(p person.Person) string { return p.Passport.IssuedOn }, func(in *person.Input, v string) { in.Passport.IssuedOn = v }),
			text("Autoridade", "passport_authority", func(p person.Person) string { return p.Passport.Authority }, func(in *person.Input, v string) { in.Passport.Authority = v }),
			text("Validade", "passport_expires_on", func(p person.Person) string { return p.Passport.ExpiresOn }, func(in *person.Input, v string) { in.Passport.ExpiresOn = v }),
		}},
	}
}

func text(p_Label, p_Name string, p_Read func(person.Person) string, p_Apply func(*person.Input, string)) formField {
	return formField{label: p_Label, name: p_Name, read: p_Read, apply: p_Apply}
}

func groupsFrom(p_Person person.Person) []ViewGroup {
	source := catalog()
	groups := make([]ViewGroup, len(source))
	for i, group := range source {
		fields := make([]ViewField, len(group.fields))
		for j, field := range group.fields {
			fields[j] = ViewField{Label: field.label, Name: field.name, Value: field.read(p_Person)}
		}
		groups[i] = ViewGroup{Title: group.title, Fields: fields}
	}
	return groups
}
