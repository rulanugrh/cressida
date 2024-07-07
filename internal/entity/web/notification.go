package web

type Notification struct {
	Content string `json:"content" form:"content"`
	Status  string `json:"status" form:"status"`
	OrderID string `json:"order_id" form:"order_id"`
	DriverName string `json:"driver_name" form:"driver_name"`
}
