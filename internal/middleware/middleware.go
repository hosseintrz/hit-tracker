package middleware

import (
	"database/sql"
	"fmt"
	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/hosseintrz/hit-tracker/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net"
)

type MiddleWare struct {
	Logger *logrus.Logger
}

func NewMiddleWare(logger *logrus.Logger) *MiddleWare {
	return &MiddleWare{
		Logger: logger,
	}
}

func (md *MiddleWare) TestMD(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		md.Logger.Info("start")
		next(c)
		md.Logger.Info("finish")
		return c.String(200, "repsonse")
	}
}

func (md *MiddleWare) HitTracker(next echo.HandlerFunc, db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		md.Logger.Info(fmt.Sprint("start md ", &c.Request().URL.Path))
		if err := next(c); err != nil {
			c.Error(err)
		}

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
			md.Logger.Error("error reading from db")
		}
		defer rows.Close()

		if !rows.Next() {
			err := crdb.Execute(func() error {
				_, err := db.Exec(`
				INSERT INTO visits (user_ip, url, count)
				VALUES ($1,$2,$3)
				 `, host, url, 1)
				if err != nil {
					md.Logger.Error("error inserting hit record: ", err.Error())
				}
				return nil
			})
			if err != nil {
				md.Logger.Error("error executing db func")
				return nil
			}
		} else {
			var visit model.Visit
			if err := rows.Scan(&visit.UserIp, &visit.URL, &visit.Count); err != nil {
				md.Logger.Error("error reading db record")
			}
			count = visit.Count + 1
			err = crdb.Execute(func() error {
				_, err := db.Exec(`
				UPDATE visits 
				SET count=count+1
				WHERE user_ip=$1 AND url=$2
			`, host, url)
				return err
			})
			if err != nil {
				md.Logger.Error("error executing db command")
				return nil
			}
		}

		md.Logger.Info(fmt.Sprintf("user with ip %s hit %s %d times", host, url, count))
		//return c.String(200, fmt.Sprintf("ip : %s hit %s %d times\n", host, url, count))
		return nil
	}
}
