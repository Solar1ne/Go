package utils

import "fmt"

// PrintHelp : helping menu
func PrintHelp() {
	fmt.Println("Доступные команды:")
	fmt.Println("help - Показать список доступных команд")
	fmt.Println("acceptOrder - Принять заказ от курьера")
	fmt.Println("returnOrder - Вернуть заказ курьеру")
	fmt.Println("deliverOrder - Выдать заказ клиенту")
	fmt.Println("getOrders - Получить список заказов")
	fmt.Println("acceptReturn - Принять возврат от клиента")
	fmt.Println("getReturns - Получить список возвратов")
	fmt.Println("pvzAdd - добавить ПВЗ")
	fmt.Println("pvzList - получить список ПВЗ")
	fmt.Println("exit - Выход из программы")
}
