package utils

import (
	"bufio"
	"strconv"
	"strings"
)

// ReadString : Reading string
func ReadString(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// ReadOrderDetails : Reading order details
func ReadOrderDetails(reader *bufio.Reader) (string, string, string) {
	orderID := ReadString(reader)
	receiverID := ReadString(reader)
	deadline := ReadString(reader)
	return orderID, receiverID, deadline
}

// ReadOrderIDs : Reading Order IDs
func ReadOrderIDs(reader *bufio.Reader) []string {
	input := ReadString(reader)
	return strings.Split(input, " ")
}

// ReadReturnDetails : Reading returning details
func ReadReturnDetails(reader *bufio.Reader) (string, string) {
	userID := ReadString(reader)
	orderID := ReadString(reader)
	return userID, orderID
}

// ReadPageDetails : Reading page details
func ReadPageDetails(reader *bufio.Reader) (int, int) {
	pageStr := ReadString(reader)
	sizeStr := ReadString(reader)
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	return page, size
}
