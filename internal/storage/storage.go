package storage

import (
	"awesomeProject3/internal/model"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// AcceptOrder : accepting order
func AcceptOrder(orderID, receiverID, deadline string) error {

	if _, err := os.Stat("orders.txt"); os.IsNotExist(err) {
		file, err := os.Create("orders.txt")
		if err != nil {
			return fmt.Errorf("ошибка создания файла: %w", err)
		}
		defer file.Close()
	}

	file, err := os.Open("orders.txt")
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == orderID {
			return fmt.Errorf("заказ с ID %s уже существует", orderID)
		}
	}

	storageDeadline, err := time.Parse("2006-01-02", deadline)
	if err != nil {
		return fmt.Errorf("неверный формат даты: %w", err)
	}
	if storageDeadline.Before(time.Now()) {
		return fmt.Errorf("срок хранения в прошлом")
	}

	file, err = os.OpenFile("orders.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла для записи: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s,%s,%s\n", orderID, receiverID, deadline))
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}
	return nil
}

// ReturnOrder : Returning order
func ReturnOrder(orderID string) error {
	if _, err := os.Stat("orders.txt"); os.IsNotExist(err) {
		return fmt.Errorf("файл с заказами не найден")
	}
	file, err := os.Open("orders.txt")
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()
	var orders []model.Order
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 4 {
			continue
		}
		order := model.Order{
			ID:              parts[0],
			ReceiverID:      parts[1],
			StorageDeadline: time.Now(),
			IsDelivered:     parts[3] == "true",
		}
		if order.ID == orderID && !order.IsDelivered && order.StorageDeadline.Before(time.Now()) {
			continue
		}
		orders = append(orders, order)
	}
	file, err = os.Create("orders.txt")
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	for _, order := range orders {
		_, err = file.WriteString(fmt.Sprintf("%s,%s,%s,%t\n", order.ID, order.ReceiverID, order.StorageDeadline.Format("2006-01-02"), order.IsDelivered))
		if err != nil {
			return fmt.Errorf("ошибка записи в файл: %w", err)
		}
	}

	fmt.Println("Заказ успешно возвращен")
	return nil
}

// DeliverOrder : returning delivered orders
func DeliverOrder(orderIDs []string) error {
	if _, err := os.Stat("orders.txt"); os.IsNotExist(err) {
		return fmt.Errorf("файл с заказами не найден")
	}
	file, err := os.Open("orders.txt")
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	var orders []model.Order
	var receiverID string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 4 {
			continue
		}
		order := model.Order{
			ID:              parts[0],
			ReceiverID:      parts[1],
			StorageDeadline: time.Now(),
			IsDelivered:     parts[3] == "true",
		}
		if contains(orderIDs, order.ID) && !order.IsDelivered && order.StorageDeadline.After(time.Now()) {
			if receiverID == order.ReceiverID {

			} else if receiverID != order.ReceiverID {
				return fmt.Errorf("все ID заказов должны принадлежать одному клиенту")
			}
			order.IsDelivered = true
			orders = append(orders, order)
		}
	}

	if len(orders) == 0 {
		return fmt.Errorf("заказы не найдены или не все заказы принадлежат одному клиенту")
	}
	file, err = os.Create("orders.txt")
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	for _, order := range orders {
		_, err = file.WriteString(fmt.Sprintf("%s,%s,%s,%t\n", order.ID, order.ReceiverID, order.StorageDeadline.Format("2006-01-02"), order.IsDelivered))
		if err != nil {
			return fmt.Errorf("ошибка записи в файл: %w", err)
		}
	}

	fmt.Println("Заказы успешно выданы")
	return nil
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// GetOrders : returning orders
func GetOrders(userID string, options ...string) ([]model.Order, error) {
	if _, err := os.Stat("orders.txt"); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл с заказами не найден")
	}
	file, err := os.Open("orders.txt")
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()
	var orders []model.Order
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 4 {
			continue
		}
		storageDeadline, err := time.Parse("2006-01-02", parts[2])
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга даты: %w", err)
		}
		order := model.Order{
			ID:              parts[0],
			ReceiverID:      parts[1],
			StorageDeadline: storageDeadline,
			IsDelivered:     parts[3] == "true",
		}

		if order.ReceiverID == userID {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

// AcceptReturn : Accepting returns
func AcceptReturn(userID, orderID string) error {
	if _, err := os.Stat("orders.txt"); os.IsNotExist(err) {
		return fmt.Errorf("файл с заказами не найден")
	}
	file, err := os.Open("orders.txt")
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	var orderToReturn *model.Order
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 5 {
			continue
		}
		storageDeadline, err := time.Parse("2006-01-02", parts[2])
		if err != nil {
			return fmt.Errorf("ошибка парсинга даты: %w", err)
		}
		deliveryDate, err := time.Parse("2006-01-022", parts[4])
		if err != nil {
			return fmt.Errorf("ошибка парсинга даты: %w", err)
		}
		order := model.Order{
			ID:              parts[0],
			ReceiverID:      parts[1],
			StorageDeadline: storageDeadline,
			IsDelivered:     parts[3] == "true",
			DeliveryDate:    deliveryDate,
		}
		if order.ID == orderID && order.ReceiverID == userID {
			orderToReturn = &order
			break
		}
	}

	if orderToReturn == nil {
		return fmt.Errorf("заказ не найден или не соответствует условиям возврата")
	}
	if time.Since(orderToReturn.DeliveryDate).Hours() > 48*24 {
		return fmt.Errorf("заказ может быть возвращен в течении 2х дней с момента выдачи")
	}
	fmt.Println("Возврат заказа успешно принят")
	return nil
}

// GetReturns : printing return list
func GetReturns(skip int, size int) {
	if _, err := os.Stat("returns.txt"); os.IsNotExist(err) {
		fmt.Println("Файл с возвратами не найден")
		return
	}
	file, err := os.Open("returns.txt")
	if err != nil {
		fmt.Printf("Ошибка открытия файла: %v\n", err)
		return
	}
	defer file.Close()

	var returns []model.Return
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 3 {
			continue
		}
		date, err := time.Parse("2006-01-02", parts[2])
		if err != nil {
			fmt.Printf("Ошибка парсинга даты: %v\n", err)
			continue
		}
		returns = append(returns, model.Return{
			OrderID: parts[0],
			UserID:  parts[1],
			Date:    date,
		})
	}
	if skip >= len(returns) {
		fmt.Println("Нет возвратов для отображения")
		return
	}
	end := skip + size
	if end > len(returns) {
		end = len(returns)
	}
	fmt.Println("Список возвратов:")
	for _, returnItem := range returns[skip:end] {
		fmt.Printf("OrderID: %s, UserID: %s, Date: %s\n", returnItem.OrderID, returnItem.UserID, returnItem.Date.Format("2006-01-02"))
	}
}

// ReadPVZDetails : Reading Details of pvz
func ReadPVZDetails(reader *bufio.Reader) (string, string, string) {
	name := readString(reader)
	address := readString(reader)
	contactInfo := readString(reader)
	return name, address, contactInfo
}

func readString(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// AddPVZ : adding pvz
func AddPVZ(name, address, contactInfo string) error {
	file, err := os.OpenFile("pvz_data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%s,%s,%s\n", name, address, contactInfo))
	return err
}

// ListPVZs : returning list of pvz
func ListPVZs() ([]PVZ, error) {
	file, err := os.Open("pvz_data.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var pvzs []PVZ
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 3 {
			continue
		}
		pvzs = append(pvzs, PVZ{Name: parts[0], Address: parts[1], ContactInfo: parts[2]})
	}

	return pvzs, nil
}

//Unused funcitions

/*func saveOrders(orders []storage.Order) error {
	file, err := os.Create("orders.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	for _, order := range orders {
		_, err := file.WriteString(fmt.Sprintf("%s,%s,%s,%t\n", order.ID, order.ReceiverID, order.StorageDeadline.Format("2006-01-02"), order.IsDelivered))
		if err != nil {
			return err
		}
	}

	return nil
}
*/
/*
func loadOrders() ([]storage.Order, error) {
	file, err := os.Open("orders.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var orders []storage.Order
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var order storage.Order
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 4 {
			continue
		}
		order.ID = parts[0]
		order.ReceiverID = parts[1]
		order.StorageDeadline, _ = time.Parse("2006-01-02", parts[2])
		order.IsDelivered = parts[3] == "true"
		orders = append(orders, order)
	}

	return orders, nil
}
*/
