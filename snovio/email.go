package snovio

// Email struct
type Email struct {
	Email     string `json:"email"`
	Status    string `json:"status"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// GetName ...
func (e Email) GetName() string {
	if len(e.FirstName) > 0 {
		if len(e.LastName) > 0 {
			return e.FirstName + " " + e.LastName
		}

		return e.FirstName
	}

	return ""
}
