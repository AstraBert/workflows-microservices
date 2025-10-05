package models

type Order struct {
	OrderId       string `json:"orderId"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Address       string `json:"address"`
	Address2      string `json:"address2"`
	City          string `json:"city"`
	State         string `json:"state"`
	Zip           string `json:"zip"`
	Country       string `json:"country"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	PaymentMethod string `json:"paymentMethod"`
	Amount        string `json:"amount"`
}

type OrdersReportStatus struct {
	OrderId string `json:"orderId"`
	Type    string `json:"type"` // payment, order, stock
	Success bool   `json:"success"`
}
