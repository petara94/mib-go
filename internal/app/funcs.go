package app

import (
	"fmt"
)

// need current and new password
func (a *App) changePassword() {
	// scan old password
	fmt.Print("Enter old password: ")
	var oldPassword string
	_, _ = fmt.Scanln(&oldPassword)

	// check old password
	if !a.currentUser.PasswordEqual(oldPassword) {
		fmt.Println("Wrong old password")
		return
	}

	// scan new password
	fmt.Print("Enter new password: ")
	var newPassword string
	_, _ = fmt.Scanln(&newPassword)

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

// updateUser update user. admin user can't change
