package entity

type Session struct {
	UserID    string `json:"userId"`
	SecretKey string `json:"secretKey"`
}
