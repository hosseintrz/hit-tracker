package main

import (
	"fmt"
	"github.com/hosseintrz/hit-tracker/internal/database"
	"github.com/hosseintrz/hit-tracker/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

func main() {

	//pageVisits := map[string]map[string]int{}
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := database.InitStore()
	if err != nil {
		log.Fatalf("couldn't open database -> \n %s \n", err.Error())
	}
	defer db.Close()

	//e.GET("/favicon.ico", func(c echo.Context) error {
	//	return nil
	//})
	e.GET("/ping", func(c echo.Context) error {
		response := struct {
			Status string
		}{"ok"}
		return c.JSON(200, response)
	})
	e.GET("/*", func(c echo.Context) error {
		return handler.RootHandler(db, c)
	})
	//mux := http.NewServeMux()
	//mux.HandleFunc("/favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
	//
	//})
	//
	//mux.HandleFunc("/ping", func(rw http.ResponseWriter, req *http.Request) {
	//	response := struct {
	//		Status string
	//	}{"ok"}
	//	res, err := json.Marshal(response)
	//	if err != nil {
	//		rw.WriteHeader(500)
	//		rw.Write([]byte("error marshaling response"))
	//		return
	//	}
	//	rw.Write(res)
	//})
	//
	//mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
	//	ip := req.RemoteAddr
	//	host, _, err := net.SplitHostPort(ip)
	//	if err != nil {
	//		rw.WriteHeader(404)
	//		rw.Write([]byte("error parsing req url"))
	//	}
	//	url := req.URL.Path
	//	if _, ok := pageVisits[host]; !ok {
	//		pageVisits[host] = map[string]int{}
	//	}
	//	hits, _ := pageVisits[host][url]
	//	pageVisits[host][url] = hits + 1
	//	fmt.Printf("%s hit with ip: %s\n", url, host)
	//	response := fmt.Sprintf("you hit %s %d times", url, pageVisits[host][url])
	//	rw.Write([]byte(response))
	//})

	defaultPort := "9090"
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = defaultPort
	}
	fmt.Printf("listening on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
