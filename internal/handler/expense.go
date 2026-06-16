package handler

import (
	"expense-tracker-api/internal/model"
	"expense-tracker-api/internal/response"
	"expense-tracker-api/internal/service"
	"fmt"
	"time"

	// "net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	service *service.ExpenseService
}

func NewExpenseHandler(s *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: s}
}

func getUserID(c *gin.Context) (int64, bool) {
	userId, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	userID := userId.(int64)
	if userID <= 0 {
		return 0, false
	}

	return userID, true
}

type CreateExpenseRequest struct {
	Title       string  `json:"title" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date" binding:"required"`
}

func (h *ExpenseHandler) Create(c *gin.Context) {
	var req CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		response.Error(c, 400, "invalid token")
		return
	}

	expense, err := h.service.Create(
		c.Request.Context(),
		userID,
		req.Title,
		req.Amount,
		req.Category,
		req.Description,
		req.Date,
	)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 201, expense)
}

func (h *ExpenseHandler) GetAll(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Error(c, 400, "invalid token")
		return
	}

	expenses, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, 500, "failed to fetch expenses")
		return
	}

	etag := fmt.Sprintf("\"%d\"", len(expenses))

	if c.GetHeader("If-None-Match") == etag {
		c.Status(304)
		return
	}

	c.Header("Cache-Control", "max-age=300")
	c.Header("ETag", etag)
	response.Success(c, 200, expenses)
}

func (h *ExpenseHandler) GetByID(c *gin.Context) {
	expenseIDStr, _ := c.Params.Get("id")
	expenseID, _ := strconv.ParseInt(expenseIDStr, 10, 64)

	userID, ok := getUserID(c)
	if !ok {
		response.Error(c, 400, "invalid token")
		return
	}

	expense, err := h.service.GetByIDAndUserID(c.Request.Context(), expenseID, userID)
	if err != nil {
		response.Error(c, 500, "failed to fetch expense")
		return
	}

	response.Success(c, 200, expense)
}

type UpdateExpenseRequest struct {
	Title       string  `json:"title" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date" binding:"required"`
}

func (h *ExpenseHandler) Update(c *gin.Context) {
	expenseIDStr, _ := c.Params.Get("id")
	expenseID, _ := strconv.ParseInt(expenseIDStr, 10, 64)

	userID, ok := getUserID(c)
	if !ok {
		response.Error(c, 400, "invalid token")
		return
	}

	if _, err := h.service.GetByIDAndUserID(c.Request.Context(), expenseID, userID); err != nil {
		response.Error(c, 400, "expense not found")
		return
	}

	var req UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	expense := model.Expense{
		ID:          expenseID,
		Title:       req.Title,
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		response.Error(c, 400, "invalid date format")
		return
	}
	expense.Date = parsedDate

	err = h.service.Update(c.Request.Context(), &expense)
	if err != nil {
		response.Error(c, 500, "failed to update expense")
		return
	}

	response.Success(c, 200, expense)
}

func (h *ExpenseHandler) Delete(c *gin.Context) {
	expenseIDStr, _ := c.Params.Get("id")
	expenseID, _ := strconv.ParseInt(expenseIDStr, 10, 64)

	userID, ok := getUserID(c)
	if !ok {
		response.Error(c, 400, "invalid token")
		return
	}

	if _, err := h.service.GetByIDAndUserID(c.Request.Context(), expenseID, userID); err != nil {
		response.Error(c, 400, "expense not found")
		return
	}

	err := h.service.Delete(c.Request.Context(), expenseID)
	if err != nil {
		response.Error(c, 500, "failed to fetch expense")
		return
	}

	response.Success(c, 200, map[string]string{"message": "expense deleted"})
}
