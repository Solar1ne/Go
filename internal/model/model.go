package model

import "time"

// Order :
type Order struct {
	ID              string
	ReceiverID      string
	StorageDeadline time.Time
	IsDelivered     bool
	DeliveryDate    time.Time
}

// Return :
type Return struct {
	OrderID string
	UserID  string
	Date    time.Time
}
