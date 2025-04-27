package auth

type User struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}