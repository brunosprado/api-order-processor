package http

import (
	"net/http"

	"time"

	"github.com/brunosprado/api-order-processor/domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateOrderRequest struct {
	Product  string `json:"product" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

func (h *handler) postOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	order := domain.Order{
		OrderID:   primitive.NewObjectID().Hex(),
		Product:   req.Product,
		Quantity:  req.Quantity,
		Status:    "CRIADO",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	if err := h.clientService.PostOrder(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist order or publish event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order_id": order.OrderID, "status": order.Status})
}
