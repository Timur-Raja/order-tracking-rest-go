package models

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type Order struct {
    ID     int    `json:"id"`
    UserID int    `json:"user_id"`
    Item   string `json:"item"`
    Amount int    `json:"amount"`
}