package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

func executeMigration(conn *pgx.Conn) {
	f, err := os.Open("db.sql")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	var b bytes.Buffer

	_, err = io.Copy(&b, f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(b.String())

	stmt, err := conn.Prepare(context.Background(), "query", b.String())
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(stmt.SQL)
	_, err = conn.Exec(context.Background(), stmt.SQL)
	if err != nil {
		log.Fatal(err)
	}
}
