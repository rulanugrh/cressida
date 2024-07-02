package web

type OrderRequest struct {
	TransporterID    uint   `json:"transporter_id" form:"transporter_id" validate:"required"`
	UserID           uint   `json:"user_id" form:"user_id" validate:"required"`
	Weight           int64  `json:"weight" form:"weight" validate:"required"`
	Distance         int64  `json:"distance" form:"distance" validate:"required"`
	Description      string `json:"description" form:"description" validate:"required"`
	PickupLat        string `json:"pickup_lat" form:"pickup_lat" validate:"required"`
	PickupLang       string `json:"pickup_lang" form:"pickup_lang" validate:"required"`
	PickupCoordinate string `json:"pickup_coordinate" form:"pickup_coordinate" validate:"required"`
	PickupAddress    string `json:"pickup_address" form:"pickup_address" validate:"required"`
	DropLat          string `json:"drop_lat" from:"drop_lat" validate:"required"`
	DropLang         string `json:"drop_lang" from:"drop_lang" validate:"required"`
	DropCoordinate   string `json:"drop_coordinate" from:"drop_coordinate" validate:"required"`
	DropAddress      string `json:"drop_address" from:"drop_address" validate:"required"`
}

type OrderResponse struct {
	UserName         string `json:"user_name" form:"user_name"`
	Weight           int64  `json:"weight" form:"weight" `
	Distance         int64  `json:"distance" form:"distance" `
	PickupLat        string `json:"pickup_lat" form:"pickup_lat"`
	PickupLang       string `json:"pickup_lang" form:"pickup_lang"`
	PickupCoordinate string `json:"pickup_coordinate" form:"pickup_coordinate"`
	PickupAddress    string `json:"pickup_address" form:"pickup_address"`
	DropLat          string `json:"drop_lat" from:"drop_lat"`
	DropLang         string `json:"drop_lang" from:"drop_lang"`
	DropCoordinate   string `json:"drop_coordinate" from:"drop_coordinate"`
	DropAddress      string `json:"drop_address" from:"drop_address"`
	Description      string `json:"description" form:"description"`
}

type TransactionResponse struct {
	TypePayment string `json:"type_payment" form:"type_payment"`
	UserEmail   string `json:"user_id" form:"user_id"`
	Weight      int64  `json:"weight" form:"weight"`
	Distance    int64  `json:"distance" form:"distance"`
	Subtotal    int64  `json:"subtotal" form:"subtotal"`
}

type UpdateOrderStatus struct {
	UUID string `json:"id" form:"id" validate:"required"`
	Status string `json:"status" form:"status" validate:"required"`
}
