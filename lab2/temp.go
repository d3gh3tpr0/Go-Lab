package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"lab2/model"
)

func main() {
	users := []model.User{
		{
			Username: "mimeo",
			Password: base64.StdEncoding.EncodeToString([]byte("hello")),
			FullName: "mimeo",
			Emails: []string{
				"mimeo@gmail.com",
			},
			Addresses: []string{
				"test",
			},
		},
	}
	data, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("users.json", data, 0644)
	if err != nil {
		panic(err)
	}
}
