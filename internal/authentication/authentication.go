package authentication

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"syscall"
)

type User struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func GetUser() (User, error) {
	home, _ := os.UserHomeDir()
	dir := fmt.Sprintf("%s/.gapp", home)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return User{}, fmt.Errorf("gapp has not been initialized, please use gapp login")
	}

	filename := fmt.Sprintf("%s/authentication.json", dir)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return User{}, err
	}

	var user User

	if err = json.Unmarshal(file, &user); err != nil {
		return User{}, err
	}

	return user, nil
}

func SaveUser() error {
	var username string

	fmt.Print("Provide your user name: ")
	fmt.Scanln(&username)

	fmt.Println("Provide your personal github token: ")
	b, _ := terminal.ReadPassword(syscall.Stdin)

	home, _ := os.UserHomeDir()

	dir := fmt.Sprintf("%s/.gapp", home)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	auth := User{
		Username: username,
		Token:    string(b),
	}

	result, err := json.Marshal(&auth)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/authentication.json", dir)
	err = ioutil.WriteFile(filename, result, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
