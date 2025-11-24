package domain

type Client struct {
	ClientID  int    `json:"client_id,omitempty" bson:"client_id,omitempty"`
	AppID     string `json:"app_id,omitempty" bson:"app_id,omitempty"`
	DomainKey string `json:"domain_key,omitempty" bson:"domain_key,omitempty"`
	Alias     string `json:"alias,omitempty" bson:"alias,omitempty"`
	Sort      string `json:"sort,omitempty" bson:"sort,omitempty"`
	Limit     int    `json:"limit,omitempty" bson:"limit,omitempty"`
}

type Order struct {
	OrderID   string `json:"order_id,omitempty" bson:"_id,omitempty"`
	Product   string `json:"product,omitempty" bson:"product,omitempty"`
	Quantity  int    `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Status    string `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt string `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type ClientService interface {
	PostOrder(order Order) error
}

type ClientStorage interface {
	PersistOrder(order Order) error
	UpdateOrderStatus(orderID string, status string) error
	GetOrderById(orderID string) (*Order, error)
}
