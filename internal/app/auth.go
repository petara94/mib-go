package app

import (
	"fmt"
	"github.com/petara94/mib-go/internal/db"
	"github.com/petara94/mib-go/internal/pkg"
)

// loginOrRegister login or register user
func (a *App) loginOrRegister() bool {
	fmt.Print("Login or register? (l/r): ")
	var answer string
	_, _ = fmt.Scanln(&answer)

	switch answer {
	case "l":
		return a.login()
	case "r":
		return a.register()
	default:
		fmt.Println("Wrong answer")
		return false
	}
}

func (a *App) login() bool {
	// scan login and pass from console
	fmt.Print("Enter login: ")
	var login string
	_, _ = fmt.Scanln(&login)

	fmt.Print("Enter password: ")
	var password string
	_, _ = fmt.Scanln(&password)

	// get user from db
	currentUser, err := a.users.GetByLogin(login)
	if err != nil {
		fmt.Println("wrong login or password")
		return false
	}

	// first run
	if currentUser.Login == "admin" && currentUser.Password == "" {
		// admin must change password
		fmt.Println("admin must change password")
		fmt.Print("Enter new password: ")
		var newPassword string
		_, _ = fmt.Scanln(&newPassword)

		currentUser.SetPassword(newPassword)

		// update currentUser
		currentUser, err = a.users.Update(login, *currentUser)
		if err != nil {
			fmt.Println("can't change password")
			return false
		}

		fmt.Println("admin password changed. please login again")

		return false
	}

	if currentUser.IsBlocked && currentUser.Login != "admin" {
		fmt.Println("user is blocked")
		return false
	}

	// check password
	if !currentUser.PasswordEqual(password) {
		fmt.Println("wrong login or password")
		return false
	}

	a.currentUser = currentUser

	return true
}

// register need login, password and password confirmation
func (a *App) register() bool {
	fmt.Print("Enter login: ")
	var login string
	_, _ = fmt.Scanln(&login)

	// check if user exists
	_, err := a.users.GetByLogin(login)
	if err == nil {

		fmt.Println("user already exists")
		return false
	}

	fmt.Print("Enter password: ")
	var password string
	_, _ = fmt.Scanln(&password)

	fmt.Print("Enter password again: ")
	var passwordAgain string
	_, _ = fmt.Scanln(&passwordAgain)

	if password != passwordAgain {
		fmt.Println("passwords are not equal")
		return false
	}

	// check password difficulty
	if !pkg.CheckPassword(password) {
		fmt.Println("password is too weak")
		return false
	}

	// read first name and last name
	fmt.Print("Enter first name: ")
	var firstName string
	_, _ = fmt.Scanln(&firstName)

	// check first name length
	if len(firstName) > 30 {
		fmt.Println("first name is too long")
		return false
	}

	fmt.Print("Enter last name: ")
	var lastName string
	_, _ = fmt.Scanln(&lastName)

	// check first name length
	if len(lastName) > 30 {
		fmt.Println("last name is too long")
		return false
	}

	// create user
	user := db.User{
		Login:     login,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		CheckPass: true,
	}
	err = a.users.Create(user)
	if err != nil {
		fmt.Println("can't create user")
		return false
	}

	fmt.Println("user created")

	// set current user
	a.currentUser = &user

	return true
}
