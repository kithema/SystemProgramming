package main

import (
	"fmt"
	"os/exec"
)

func main() {
    fmt.Println("start | kill | exit")
    input := ""
    for{
        fmt.Println("Enter commad...")
        fmt.Scan(&input)
		//свитч с выбром команды
        switch input {
		case "exit":
			fmt.Println("Exit...")
			return

		case "start":
            fmt.Println("Enter name: ")
            fmt.Scan(&input)
			//запуск желаемого приложения 
            cmd := exec.Command(input)
            cmd.Start()

		case "kill":
            fmt.Println("Enter name: ")
            fmt.Scan(&input)
            killProgram(input)

		default:
			fmt.Printf("Error: %s\n", input)
		}
    }

}
//возвращает error но я забыл добавить его обработку
func killProgram(name string) error {
  //Закрытие потока через taskkill
  cmd := exec.Command("taskkill", "/IM", name, "/F")
  err := cmd.Run()
  if err != nil {
    return fmt.Errorf("taskkill failed: %v", err)
  }
  fmt.Printf("Process %s killed successfully\n", name)
  return nil

}
