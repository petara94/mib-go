package app

import (
	"fmt"
)

const (
	GetUsersWithPagination = "get users with pagination"
	GetUser                = "get user"
	UpdateUser             = "update user"
	CreateUser             = "create user"
	CreateUserWithLogin    = "create user with login"
	DeleteUser             = "delete user"

	SetBlock = "set block user"
	SetCheck = "set check password"

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
		GetUsersWithPagination,
		GetUser,
		UpdateUser,
		CreateUser,
		CreateUserWithLogin,
		DeleteUser,
		SetBlock,
		SetCheck,
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
