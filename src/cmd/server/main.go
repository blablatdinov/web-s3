package main

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber"
)

func healthCheck(fiberContext *fiber.Ctx) error {
	return fiberContext.JSON(fiber.Map{
		"app":      true,
		"postgres": true,
		"redis":    true,
		"s3":       true,
	})
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.ListenAndServe(":8090", nil)
}
