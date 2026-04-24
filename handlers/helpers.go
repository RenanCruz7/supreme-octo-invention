package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func parsePagination(c *gin.Context) (page, limit int) {
	page = 1
	limit = 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	return page, limit
}
