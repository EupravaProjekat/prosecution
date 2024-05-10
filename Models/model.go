package Models

type User struct {
	Email    string
	Role     string
	Requests []Request
}

// SpecialCross or temporaryExport
type Request struct {
	uuid           string
	RequestType    string
	CarPlateNumber string
	Description    string
}
