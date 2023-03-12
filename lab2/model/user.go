package model

type User struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	FullName  string   `json:"fullName"`
	Emails    []string `json:"emails"`
	Addresses []string `json:"addresses"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
