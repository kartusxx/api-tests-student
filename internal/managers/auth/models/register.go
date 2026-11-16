package models

// описывает тело запроса регистрации пользователя
type RegisterRequest struct {
	BirthDate string `json:"birth_date"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}
