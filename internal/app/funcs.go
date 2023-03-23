package app

import (
	"encoding/json"
	"fmt"
	"github.com/petara94/mib-go/internal/db"
	"github.com/petara94/mib-go/internal/pkg"
)

// print json user
func printUser(user db.User) {
	data, err := json.MarshalIndent(user, "", "  ")

	if err != nil {
		err = fmt.Errorf("Can't print user:", err)
		panic(err)
	}

	fmt.Println(string(data))
}

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

// getUsersWithPagination returns users with pagination
func (a *App) getUsersWithPagination() {
	// scan page
	fmt.Print("Enter page: ")
	var page int
	_, err := fmt.Scanln(&page)
	if err != nil {
		fmt.Println("Wrong page value")
		return
	}

	// scan limit
	fmt.Print("Enter limit: ")
	var limit int
	_, err = fmt.Scanln(&limit)
	if err != nil {
		fmt.Println("Wrong limit value")
		return
	}

	// get users
	users, err := a.users.GetWithPagination(page, limit)
	if err != nil {
		fmt.Println("Can't get users:", err)
		return
	}

	// print users
	for _, user := range users {
		printUser(*user)
	}
}

// getUser returns user by login
func (a *App) getUser() {
	// scan login
	fmt.Print("Enter login: ")
	var login string
	_, _ = fmt.Scanln(&login)

	// get user
	user, err := a.users.GetByLogin(login)
	if err != nil {
		fmt.Println("Can't get user:", err)
		return
	}

	printUser(*user)
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

// updateUser updates all user fields. admin can't change login and block status
func (a *App) updateUser() {
	// scan login
	fmt.Print("Enter login: ")
	var login string
	_, _ = fmt.Scanln(&login)

	// get user
	user, err := a.users.GetByLogin(login)
	if err != nil {
		fmt.Println("Can't get user:", err)
		return
	}

	// admin can't change login and block status
	if login != "admin" {
		// scan new login
		fmt.Print("Enter new login: ")
		_, err = fmt.Scanln(&user.Login)
		if err != nil {
			fmt.Println("Wrong login value")
			return
		}
	}

	// scan password
	password := pkg.ReadPassword("Enter password: ")

	// scan check password
	fmt.Print("Enter check password: ")
	var checkPassword bool
	_, err = fmt.Scanln(&checkPassword)
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

	// admin can't change login and block status
	if login != "admin" {
		// scan block status
		fmt.Print("Enter block status: ")
		_, err = fmt.Scanln(&user.IsBlocked)
		if err != nil {
			fmt.Println("Wrong block status value")
			return
		}
	}

	// update user
	user.SetPassword(password)
	user.CheckPass = checkPassword
	user.FirstName = firstName
	user.LastName = lastName

	user, err = a.users.Update(login, *user)
	if err != nil {
		fmt.Println("Can't update user:", err)
		return
	}

	fmt.Println("User updated:")
	printUser(*user)
}

// deleteUser deletes user by login
func (a *App) deleteUser() {
	// scan login
	fmt.Print("Enter login: ")
	var login string
	_, _ = fmt.Scanln(&login)

	// check login
	if login == "admin" {
		fmt.Println("Can't delete admin")
		return
	}

	// delete user
	err := a.users.Delete(login)
	if err != nil {
		fmt.Println("Can't delete user:", err)
		return
	}

	fmt.Println("User deleted")
}

// setCheckPassword sets check password flag
func (a *App) setCheckPassword() {
	// scan login
	fmt.Print("Enter login: ")
	var login string
	_, _ = fmt.Scanln(&login)

	// get user
	user, err := a.users.GetByLogin(login)
	if err != nil {
		fmt.Println("Can't get user:", err)
		return
	}

	// scan check password
	fmt.Print("Enter check password: ")
	var checkPassword bool
	_, err = fmt.Scanln(&checkPassword)
	if err != nil {
		fmt.Println("Wrong check password value")
		return
	}

	// update user
	user.CheckPass = checkPassword
	user, err = a.users.Update(login, *user)
	if err != nil {
		fmt.Println("Can't update user:", err)
		return
	}

	fmt.Println("Check password updated")
	printUser(*user)
}

// setBlock sets block status. admin can't be blocked
func (a *App) setBlock() {
	// scan login
	fmt.Print("Enter login: ")
	var login string
	_, _ = fmt.Scanln(&login)

	// admin can't be blocked
	if login == "admin" {
		fmt.Println("Can't block admin")
		return
	}

	// get user
	user, err := a.users.GetByLogin(login)
	if err != nil {
		fmt.Println("Can't get user:", err)
		return
	}

	// scan block status
	fmt.Print("Enter block status: ")
	_, err = fmt.Scanln(&user.IsBlocked)
	if err != nil {
		fmt.Println("Wrong block status value")
		return
	}

	// update user
	user, err = a.users.Update(login, *user)
	if err != nil {
		fmt.Println("Can't update user:", err)
		return
	}

	fmt.Println("Block status updated")
	printUser(*user)
}

// createWithLogin creates user with login
func (a *App) createWithLogin() {
	// scan login
	fmt.Print("Enter login: ")
	var login string
	_, _ = fmt.Scanln(&login)

	user := db.User{Login: login}

	err := a.users.Create(user)
	if err != nil {
		fmt.Println("Can't create user:", err)
		return
	}

	fmt.Println("User created:")
	printUser(user)
}
