package handler

import (
	"database/sql"
	"fmt"
	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/hosseintrz/hit-tracker/internal/model"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
)

func RootHandler(db *sql.DB, c echo.Context) error {
	count := 1
	ip := c.Request().RemoteAddr
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		return c.String(400, "error parsing req url")
	}
	url := c.Request().URL.Path
	rows, err := db.Query(`
		SELECT * FROM visits
			WHERE user_ip=$1 and url=$2
		`, host, url)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		err := crdb.Execute(func() error {
			_, err := db.Exec(`
				INSERT INTO visits (user_ip, url, count)
				VALUES ($1,$2,$3)
				 `, host, url, 1)

			if err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}
			return nil
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	} else {
		var visit model.Visit
		if err := rows.Scan(&visit.UserIp, &visit.URL, &visit.Count); err != nil {
			return err
		}
		count = visit.Count + 1
		err := crdb.Execute(func() error {
			_, err := db.Exec(`
				UPDATE visits 
				SET count=count+1
				WHERE user_ip=$1 AND url=$2
			`, host, url)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}
			return nil
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}
	c.Logger().Info(fmt.Sprintf("user with ip %s hit %s %d times", host, url, count))
	return c.String(200, fmt.Sprintf("ip : %s hit %s %d times\n", host, url, count))
}
