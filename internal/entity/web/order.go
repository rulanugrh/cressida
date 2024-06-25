package web

type OrderRequest struct {
	TransporterID    uint   `json:"transporter_id" form:"transporter_id" validate:"required"`
	UserID           uint   `json:"user_id" form:"user_id" validate:"required"`
	Weight           int64  `json:"weight" form:"weight" validate:"required"`
	Distance         int64  `json:"distance" form:"distance" validate:"required"`
	PickupLocation   string `json:"pickup_location" form:"pickup_location" validate:"required"`
	DeliveryLocation string `json:"delivery_location" form:"delivery_location" validate:"required"`
	Description      string `json:"description" form:"description" validate:"required"`
	Status           string `json:"status" form:"status" validate:"required"`
}

type OrderResponse struct {
	UserName         string `json:"user_name" form:"user_name"`
	Weight           int64  `json:"weight" form:"weight" `
	Distance         int64  `json:"distance" form:"distance" `
	PickupLocation   string `json:"pickup_location" form:"pickup_location"`
	DeliveryLocation string `json:"delivery_location" form:"delivery_location" `
	Description      string `json:"description" form:"description" `
}

type TransactionResponse struct {
	TypePayment string `json:"type_payment" form:"type_payment"`
	UserEmail   string `json:"user_id" form:"user_id"`
	Weight      int64  `json:"weight" form:"weight"`
	Distance    int64  `json:"distance" form:"distance"`
	Subtotal    int64  `json:"subtotal" form:"subtotal"`
}
