package app

import (
	"fmt"
	"github.com/petara94/mib-go/internal/db"
	"github.com/petara94/mib-go/internal/services"
)

type App struct {
	users       services.UserRepo
	currentUser *db.User
}

func NewApp(users services.UserRepo) *App {
	return &App{users: users}
}

func (a *App) Run() error {
	for i := 0; !a.loginOrRegister(); i++ {
		if i == 2 {
			fmt.Println("too many attempts")
			return nil
		}
	}

	switch a.currentUser.Login {
	case "admin":
		a.adminLoop()
	default:
		a.userLoop()
	}

	return nil
}

func (a *App) adminLoop() {
	for {
		PrintCommands(AdminCommands)
		command := GetAdminCommand(a.getIntFromConsole())

		switch command {
		case "create_user":
		case "change_password":
			a.changePassword()
		case "delete_user":
		case "set_block_user":
		case "set_check_password":
		case "set_admin":
		case "exit":
			fmt.Println("bye")
			return
		default:
			fmt.Println("unknown command")
			return
		}
	}
}

func (a *App) userLoop() {
	for {
		PrintCommands(UserCommands)
		command := GetUserCommand(a.getIntFromConsole())

		switch command {
		case "change_password":
			a.changePassword()
		case "exit":
			fmt.Println("bye")
			return
		default:
			fmt.Println("unknown command")
			return
		}
	}
}

func (a *App) getIntFromConsole() int {
	fmt.Printf("%s@db $ ", a.currentUser.Login)
	var commandIndex int
	_, err := fmt.Scanln(&commandIndex)
	if err != nil {
		return -1
	}
	return commandIndex
}
