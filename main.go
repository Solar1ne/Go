package main

import (
	"bufio"
	"fmt"
	"gitlab.ozon.dev/vtikunov/go-11-junior-project/internal/storage"
	"gitlab.ozon.dev/vtikunov/go-11-junior-project/internal/utils"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

var mutex = &sync.Mutex{}

func main() {
	runtime.GOMAXPROCS(0)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Введите команду:")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	for {
		fmt.Println("Приложение работает...")
		time.Sleep(10 * time.Millisecond)
		fmt.Print("> ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "help":
			utils.PrintHelp()
		case "acceptOrder":
			fmt.Println("Введите ID заказа, ID получателя и срок хранения:")
			orderID, receiverID, deadline := utils.ReadOrderDetails(reader)
			go func() {
				mutex.Lock()
				defer mutex.Unlock()
				err := storage.AcceptOrder(orderID, receiverID, deadline)
				if err != nil {
					fmt.Println("Ошибка:", err)
				} else {
					fmt.Println("Заказ успешно принят")
				}
			}()
		case "returnOrder":
			fmt.Println("Введите ID заказа:")
			orderID := utils.ReadString(reader)
			go func() {
				mutex.Lock()
				defer mutex.Unlock()
				err := storage.ReturnOrder(orderID)
				if err != nil {
					fmt.Println("Ошибка:", err)
				} else {
					fmt.Println("Заказ успешно возвращен")
				}
			}()
		case "deliverOrder":
			fmt.Println("Введите ID заказов через пробел:")
			orderIDs := utils.ReadOrderIDs(reader)
			go func() {
				mutex.Lock()
				defer mutex.Unlock()
				err := storage.DeliverOrder(orderIDs)
				if err != nil {
					fmt.Println("Ошибка:", err)
				} else {
					fmt.Println("Заказы успешно выданы")
				}
			}()
		case "getOrders":
			fmt.Println("Введите ID пользователя:")
			userID := utils.ReadString(reader)
			go func() {
				mutex.Lock()
				defer mutex.Unlock()
				orders, err := storage.GetOrders(userID)
				if err != nil {
					fmt.Println("Ошибка:", err)
				} else {
					fmt.Println("Список заказов:")
					for _, order := range orders {
						fmt.Printf("ID: %s, Получатель: %s, Срок хранения: %s\n", order.ID, order.ReceiverID, order.StorageDeadline.Format("2006-01-02"))
					}
				}
			}()
		case "acceptReturn":
			fmt.Println("Введите ID пользователя и ID заказа:")
			userID, orderID := utils.ReadReturnDetails(reader)
			go func() {
				mutex.Lock()
				defer mutex.Unlock()
				err := storage.AcceptReturn(userID, orderID)
				if err != nil {
					fmt.Println("Ошибка:", err)
				} else {
					fmt.Println("Возврат успешно принят")
				}
			}()
		case "pvzAdd":
			fmt.Println("Введите название, адрес и контактные данные ПВЗ:")
			name, address, contactInfo := storage.ReadPVZDetails(reader)
			go func() {
				mutex.Lock()
				defer mutex.Unlock()
				err := storage.AddPVZ(name, address, contactInfo)
				if err != nil {
					fmt.Println("Ошибка:", err)
				} else {
					fmt.Println("ПВЗ успешно добавлен")
				}
			}()
		case "pvzList":
			go func() {
				mutex.Lock()
				defer mutex.Unlock()
				pvzs, err := storage.ListPVZs()
				if err != nil {
					fmt.Println("Ошибка:", err)
				} else {
					fmt.Println("Список ПВЗ:")
					for _, pvz := range pvzs {
						fmt.Printf("Название: %s, Адрес: %s, Контактные данные: %s\n", pvz.Name, pvz.Address, pvz.ContactInfo)
					}
				}
			}()
		case "getReturns":
			fmt.Println("Введите номер страницы и количество элементов на странице:")
			page, size := utils.ReadPageDetails(reader)
			go func() {
				mutex.Lock()
				defer mutex.Unlock()
				storage.GetReturns(page, size)
			}()
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("Неизвестная команда. Введите 'help' для получения списка доступных команд.")
		}
	}
}
