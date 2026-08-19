package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	blog, err := NewBlog()
	if err != nil {
		log.Printf("initialize blog: %v", err)
		os.Exit(1)
	}

	sqlDB, err := blog.DB.DB()
	if err != nil {
		log.Printf("get sql database: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	server := &http.Server{Addr: ":8086", Handler: blog.Router()}
	log.Printf("blog api listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("serve blog api: %v", err)
		os.Exit(1)
	}
}
