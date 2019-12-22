package addressparser

type ParsedAddressPart struct {
	ShortName  string
	Name       string
	Id         string
	ParentId   string
	Level      int
	PostalCode string
}
