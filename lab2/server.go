package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lab2/model"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func handleConnection(conn net.Conn, db map[string]*model.User) {
	defer conn.Close()
	// generate random number
	rand.Seed(time.Now().Unix())
	result := rand.Intn(100) + 1
	// send welcome message to client
	welcome := "Welcome to the guessing game! Guess a number between 1 and 100. Please login first!\nUsername:\t"
	_, err := conn.Write([]byte(welcome))
	if err != nil {
		fmt.Println("error sending data: ", err)
		return

	}
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println("error reading data: ", err)
		return
	}
	username := strings.TrimSpace(string(data[:n]))
	user := db[username]
	if user == nil {
		_, err := conn.Write([]byte("User not found"))
		if err != nil {
			fmt.Println("error sending data ", err)
			return
		}
		return
	}
	_, err = conn.Write([]byte("Password:\t"))
	if err != nil {
		fmt.Println("error sending data: ", err)
		return
	}
	data = make([]byte, 1024)
	n, err = conn.Read(data)
	password := strings.TrimSpace(string(data[:n]))
	if !authentication(password, user) {
		_, err = conn.Write([]byte("Password is incorrect"))
		if err != nil {
			fmt.Println("error sending data: ", err)
			return
		}
		return
	}
	_, err = conn.Write([]byte("Login successfully!"))
	if err != nil {
		fmt.Println("error sending data ", err)
		return
	}
	// process guesses from client
	for {
		// read data from client
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println("error reading data: ", err)
			return
		}
		// parse guess from client
		guessString := string(data[:n])
		if strings.TrimSpace(guessString) == "quit" {
			fmt.Println("client disconnected")
			return
		}
		guess, err := strconv.Atoi(strings.TrimSpace(guessString))
		if err != nil {
			// invalid guess
			_, err = conn.Write([]byte("Invalid guess. Please enter a number between 1 and 100."))
			if err != nil {
				fmt.Println("error sending data: ", err)
				return
			}
			continue
		}
		// check guess against result
		if guess < result {
			_, err = conn.Write([]byte("Too small."))
			if err != nil {
				fmt.Println("error sending data: ", err)
				return
			}
		} else if guess > result {
			_, err = conn.Write([]byte("Too large."))
			if err != nil {
				fmt.Println("error sending data: ", err)
				return
			}
		} else {
			// correct guess
			_, err = conn.Write([]byte("Congratulations, you guessed the number!"))
			if err != nil {
				fmt.Println("error sending data: ", err)
				return
			}
			// ask if client wants to play again
			_, err = conn.Write([]byte("Play again? (y/n)"))
			if err != nil {
				fmt.Println("error sending data: ", err)
				return
			}
			// read response from client
			data = make([]byte, 1024)
			n, err = conn.Read(data)
			if err != nil {
				fmt.Println("error reading data: ", err)
				return
			}
			response := strings.TrimSpace(string(data[:n]))
			if response == "y" {
				// generate new result and start new game
				result = rand.Intn(100) + 1
				_, err = conn.Write([]byte("Welcome to the guessing game! Guess a number between 1 and 100."))
				if err != nil {
					fmt.Println("error sending data: ", err)
					return
				}
			} else {
				// end game
				fmt.Println("client disconnected")
				return
			}
		}
	}
}

func main() {
	db := getUsers()
	// listen on TCP port 1234
	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("error listening: ", err)
		return
	}
	defer ln.Close()
	fmt.Println("listening on port 1234")
	// accept
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("error accepting connection: ", err)
			continue
		}
		fmt.Println("client connected: ", conn.RemoteAddr())
		go handleConnection(conn, db)
	}
}

func getUsers() map[string]*model.User {
	data, err := ioutil.ReadFile("users.json")
	if err != nil {
		panic(err)
	}
	var users []model.User
	err = json.Unmarshal(data, &users)
	if err != nil {
		panic(err)
	}
	ans := make(map[string]*model.User)
	for _, user := range users {
		ans[user.Username] = &user
	}
	return ans
}

func authentication(password string, user *model.User) bool {
	decodedPassword, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		panic(err)
	}
	passwordTmp := string(decodedPassword)
	if passwordTmp == password {
		return true
	}
	return false
}
