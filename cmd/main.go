package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dasilro/gocourse_course/internal/course"
	"github.com/dasilro/gocourse_course/pkg/bootstrap"
	"github.com/dasilro/gocourse_course/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	l := bootstrap.InitLogger()

	db, err := bootstrap.DBConnection()
	if err != nil {
		l.Fatal(err)
	}

	ctx := context.Background()
	courseRepo := course.NewRepo(l, db)
	courseSrv := course.NewService(l, courseRepo)

	h := handler.NewCourseHTTPServer(ctx, course.MakeEndpoints(courseSrv, os.Getenv("PAGINATOR_LIMIT_DEFAULT")))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)
	srv := &http.Server{
		Handler:      accessControl(h),
		Addr:         address,
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
	}

	errCh := make(chan error)
	go func() {
		l.Println("listen in ", address)
		errCh <- srv.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		log.Fatal(err)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, PUT, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type")

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})

}
