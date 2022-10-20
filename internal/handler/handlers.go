package handler

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net"
)

type Handler struct {
	logger *logrus.Logger
}

func NewHandler(logger *logrus.Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) LoginHandler(c echo.Context) error {
	return c.String(200, "welcome to login page")
}

func (h *Handler) GetStats(c echo.Context, db *sql.DB) error {
	ip := c.Request().RemoteAddr
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		c.Error(err)
	}
	rows, err := db.Query(`
		SELECT * FROM visits
		WHERE user_ip=$1
	`, host)
	if err != nil {
		h.handleErr(err, "errro getting state", c)
	}
	type Record struct {
		URL   string `json:"url"`
		Count int    `json:"count"`
	}
	records := []Record{}
	for rows.Next() {
		var userIp, url string
		var cnt int
		err := rows.Scan(&userIp, &url, &cnt)
		if err != nil {
			h.handleErr(err, "error scanning records", c)
		}
		records = append(records, Record{URL: url, Count: cnt})
	}
	return c.JSON(200, records)
}

func (h *Handler) PaymentHandler(c echo.Context) error {
	return c.String(200, "payment pending")
}

func (h *Handler) handleErr(err error, msg string, c echo.Context) {
	h.logger.Error(fmt.Sprint(msg, " -> ", err.Error()))
	c.Error(err)
}
