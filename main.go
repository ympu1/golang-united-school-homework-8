package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func Perform(args Arguments, writer io.Writer) error {
	switch args["operation"] {
	case "add":
		return add(args, writer)
	case "list":
		return list(args, writer)
	case "findById":
		return findById(args, writer)
	case "remove":
		return remove(args, writer)
	case "":
		return fmt.Errorf(operationError)
	default:
		return fmt.Errorf(operationNotAllowed, args["operation"])
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() Arguments {
	var args Arguments

	return args
}

func remove(args Arguments, writer io.Writer) error {
	fileName := args["fileName"]
	if fileName == "" {
		return fmt.Errorf(fileNameError)
	}

	id := args["id"]
	if id == "" {
		return fmt.Errorf(idError)
	}

	users, err := readUsersFromFile(fileName)
	if err != nil {
		return err
	}

	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)

			return writeUsersToFile(users, fileName)
		}
	}

	errorText := fmt.Sprintf(userIdNotFound, id)
	_, err = writer.Write([]byte(errorText))
	return err
}

func findById(args Arguments, writer io.Writer) error {
	fileName := args["fileName"]
	if fileName == "" {
		return fmt.Errorf(fileNameError)
	}

	id := args["id"]
	if id == "" {
		return fmt.Errorf(idError)
	}

	users, err := readUsersFromFile(fileName)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Id == id {
			userBytes, err := json.Marshal(user)
			if err != nil {
				return err
			}

			_, err = writer.Write(userBytes)
			return err
		}
	}

	return nil
}

func add(args Arguments, writer io.Writer) error {
	fileName := args["fileName"]
	if fileName == "" {
		return fmt.Errorf(fileNameError)
	}

	item := args["item"]
	if item == "" {
		return fmt.Errorf(itemError)
	}

	var newUser User
	err := json.Unmarshal([]byte(item), &newUser)
	if err != nil {
		return err
	}

	users, err := readUsersFromFile(fileName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	for _, user := range users {
		if user.Id == newUser.Id {
			errorText := fmt.Sprintf(userIdExist, newUser.Id)
			_, err = writer.Write([]byte(errorText))
			return err
		}
	}

	users = append(users, newUser)
	return writeUsersToFile(users, fileName)
}

func writeUsersToFile(users []User, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePermissionConst)
	defer file.Close()
	if err != nil {
		return err
	}

	usersBytes, err := json.Marshal(users)
	if err != nil {
		return err
	}

	_, err = file.Write(usersBytes)
	return err
}

func readUsersFromFile(fileName string) ([]User, error) {
	var users []User

	file, err := os.OpenFile(fileName, os.O_RDONLY, filePermissionConst)
	defer file.Close()
	if err != nil {
		return users, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return users, err
	}

	err = json.Unmarshal(bytes, &users)
	return users, err
}

func list(args Arguments, writer io.Writer) error {
	fileName := args["fileName"]

	if fileName == "" {
		return fmt.Errorf(fileNameError)
	}

	users, err := readUsersFromFile(fileName)
	if err != nil {
		return err
	}

	usersBytes, err := json.Marshal(users)
	if err != nil {
		return err
	}
	_, err = writer.Write(usersBytes)
	if err != nil {
		return err
	}

	return nil
}
