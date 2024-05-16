package Models

type BreachType int

const (
    BreachTypeUnknown BreachType = iota
    BreachTypeSpeeding
    BreachTypeParkingViolation
    BreachTypeTrafficSignalViolation
    // Add more breach types as needed
)

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

type Prosecution struct {
	JMBG         string
	TypeOfBreach BreachType
}

