package main

const filePermissionConst = 0644
const fileNameError = "-fileName flag has to be specified"
const itemError = "-item flag has to be specified"
const operationNotAllowed = "Operation %s not allowed!"
const operationError = "-operation flag has to be specified"
const userIdExist = "Item with id %s already exists"
const userIdNotFound = "Item with id %s not found"
const idError = "-id flag has to be specified"

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
