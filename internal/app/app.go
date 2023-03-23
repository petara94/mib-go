package app

import (
	"fmt"
	"github.com/petara94/mib-go/internal/db"
	"github.com/petara94/mib-go/internal/services"
	"strconv"
	"time"
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
		time.Sleep(time.Second)
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
		case CreateUser:
			a.createUser()
		case ChangePassword:
			a.changePassword()
		case DeleteUser:
		case SetCheck:
		case SetBlock:
		case SetAdmin:
		case Exit:
			fmt.Println("bye")
			return
		default:
			fmt.Println("unknown command")
		}
		time.Sleep(time.Second)
	}
}

func (a *App) userLoop() {
	for {
		PrintCommands(UserCommands)
		command := GetUserCommand(a.getIntFromConsole())

		switch command {
		case ChangePassword:
			a.changePassword()
		case Exit:
			fmt.Println("bye")
			return
		default:
			fmt.Println("unknown command")
		}
		time.Sleep(time.Second)
	}
}

func (a *App) getIntFromConsole() int {
	var (
		commandIndex int
		err          error
	)

	for {
		fmt.Printf("%s@db $ ", a.currentUser.Login)
		var command string
		_, _ = fmt.Scanln(&command)

		if command == "" {
			continue
		}

		// convert string to int
		commandIndex, err = strconv.Atoi(command)
		if err != nil {
			return -1
		}

		return commandIndex
	}
}
