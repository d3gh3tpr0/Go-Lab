package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// connect to server on TCP port 1234
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("error connecting: ", err)
		return
	}
	defer conn.Close()
	// read welcome message from server
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println("error reading data: ", err)
		return
	}
	fmt.Print(string(data[:n]))
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	_, err = conn.Write([]byte(username))
	if err != nil {
		fmt.Println("error sending data ", err)
		return
	}
	data = make([]byte, 1024)
	n, err = conn.Read(data)
	if err != nil {
		fmt.Println("error reading data: ", err)
		return
	}
	response := strings.TrimSpace(string(data[:n]))
	if strings.Contains(response, "Password") {
		fmt.Print(response + "\t")
		password, _ := reader.ReadString('\n')
		_, err = conn.Write([]byte(password))
		if err != nil {
			fmt.Println("error sending data ", err)
			return
		}
		fmt.Println("here 1")
		data = make([]byte, 1024)
		n, err = conn.Read(data)
		response = strings.TrimSpace(string(data[:n]))
		if strings.Contains(response, "Password is incorrect") {
			fmt.Println(response)
			return
		}
	} else {
		fmt.Println(response)
		return
	}
	// process guesses and responses from server
	for {
		// read guess from user

		fmt.Print("Guess: ")
		guessString, _ := reader.ReadString('\n')
		// send guess to server
		_, err = conn.Write([]byte(guessString))
		if err != nil {
			fmt.Println("error sending data: ", err)
			return
		}
		// read response from server
		data = make([]byte, 1024)
		n, err = conn.Read(data)
		if err != nil {
			fmt.Println("error reading data: ", err)
			return
		}
		response := string(data[:n])
		fmt.Println(response)
		// check if game is over
		if strings.Contains(response, "Congratulations") {
			// ask if user wants to play again
			reader = bufio.NewReader(os.Stdin)
			fmt.Print("Play again? (y/n): ")
			responseString, _ := reader.ReadString('\n')
			// send response to server
			_, err = conn.Write([]byte(responseString))
			if err != nil {
				fmt.Println("error sending data: ", err)
				return
			}
			// read welcome message from server or end game
			data = make([]byte, 1024)
			n, err = conn.Read(data)
			if err != nil {
				fmt.Println("error reading data: ", err)
				return
			}
			if strings.Contains(string(data[:n]), "Welcome") {
				fmt.Println(string(data[:n]))
			} else {
				fmt.Println("Goodbye!")
				return
			}
		}
	}
}
