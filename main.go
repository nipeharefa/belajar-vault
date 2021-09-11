package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	vaultgo "github.com/mittwald/vaultgo"
)

type DBCred struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Response struct {
	Data DBCred `json:"data"`
}

func main() {

	c, err := vaultgo.NewClient("http://localhost:8200", nil, vaultgo.WithAuthToken("s.FK8GzbHDgP4GhK7jjhJEXopQ"))
	if err != nil {
		log.Fatal(err)
		return
	}
	// c.Token()
	fmt.Println(c, c.Token())

	var a Response
	// c.Secr
	err = c.Read([]string{"v1/postgres/creds/nginx"}, &a, nil)
	fmt.Println(a, err)

	urlExample := fmt.Sprintf("postgres://%s:%s@localhost:15432/postgres", a.Data.Username, a.Data.Password)
	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// executeMigration(conn)

	fmt.Println(conn.Ping(context.Background()))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		stmt, err := conn.Prepare(r.Context(), "get-config", "SELECT * FROM roles")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		rows, err := conn.Query(r.Context(), stmt.SQL)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		defer rows.Close()

		w.Write([]byte(stmt.Name))
	})
	http.ListenAndServe(":3000", r)

}
