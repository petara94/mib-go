package app

import (
	"encoding/json"
	"fmt"
	"github.com/petara94/mib-go/internal/db"
	"github.com/petara94/mib-go/internal/pkg"
)

// need current and new password
func (a *App) changePassword() {
	// scan old password
	oldPassword := pkg.ReadPassword("Enter old password: ")

	// check old password
	if !a.currentUser.PasswordEqual(oldPassword) {
		fmt.Println("Wrong old password")
		return
	}

	// scan new password
	newPassword := pkg.ReadPassword("Enter new password: ")

	if newPassword == "" || newPassword == oldPassword {
		fmt.Println("Wrong new password")
		return
	}

	// update password
	a.currentUser.SetPassword(newPassword)

	// update user
	var err error
	a.currentUser, err = a.users.Update(a.currentUser.Login, *a.currentUser)
	if err != nil {
		fmt.Println("Can't change password:", err)
		return
	}

	fmt.Println("Password changed")
}

// createUser creates new user
func (a *App) createUser() {
	// scan login
	fmt.Print("Enter login: ")
	var login string
	_, _ = fmt.Scanln(&login)

	// check login
	if _, err := a.users.GetByLogin(login); err == nil {
		fmt.Println("User already exists")
		return
	}

	// scan password
	password := pkg.ReadPassword("Enter password: ")

	// scan check password
	fmt.Print("Enter check password: ")
	var checkPassword bool
	_, err := fmt.Scanln(&checkPassword)
	if err != nil {
		fmt.Println("Wrong check password value")
		return
	}

	if checkPassword && !pkg.CheckPassword(password) {
		fmt.Println("password is not strong enough")
		return
	}

	// scan first name
	fmt.Print("Enter first name: ")
	var firstName string
	_, _ = fmt.Scanln(&firstName)

	// scan last name
	fmt.Print("Enter last name: ")
	var lastName string
	_, _ = fmt.Scanln(&lastName)

	user := db.User{
		Login:     login,
		FirstName: firstName,
		LastName:  lastName,
		CheckPass: checkPassword,
	}
	user.SetPassword(password)

	err = a.users.Create(user)
	if err != nil {
		fmt.Println("Can't create user:", err)
		return
	}

	fmt.Println("User created:")
	printUser(user)
}

// print json user
func printUser(user db.User) {
	data, err := json.MarshalIndent(user, "", "  ")

	if err != nil {
		err = fmt.Errorf("Can't print user:", err)
		panic(err)
	}

	fmt.Println(string(data))
}
