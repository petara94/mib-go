package app

import (
	"fmt"
)

const (
	GetAllUsers = "get all users"
	GetUser     = "get user"
	UpdateUser  = "update user"
	CreateUser  = "create user"
	DeleteUser  = "delete user"

	SetBlock = "set block user"
	SetCheck = "set check password"
	SetAdmin = "set admin"

	ChangePassword = "change password"

	Exit = "exit"
)

var (
	UserCommands = []string{
		Exit,
		ChangePassword,
	}

	AdminCommands = []string{
		Exit,
		ChangePassword,
		GetAllUsers,
		GetUser,
		UpdateUser,
		CreateUser,
		DeleteUser,
		SetBlock,
		SetCheck,
		SetAdmin,
	}
)

// PrintCommands print commands
func PrintCommands(commands []string) {
	fmt.Println("choose command:")
	for i, command := range commands {
		fmt.Printf(" %d) %s\n", i, command)
	}
	fmt.Println()
}

// GetUserCommand get users command by index
func GetUserCommand(index int) string {
	if index < 0 || index >= len(UserCommands) {
		return ""
	}

	return UserCommands[index]
}

// GetAdminCommand get admins command by index
func GetAdminCommand(index int) string {
	if index < 0 || index >= len(AdminCommands) {
		return ""
	}

	return AdminCommands[index]
}
