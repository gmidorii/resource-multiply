package resourcemultiply

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func MultiplySchema(schema string, copyNum int) (error, func() error) {
	cmd := exec.Command("task", "dump-db") // TODO: pg_dump へ変更
	cmd.Run()

	f, err := os.Open("./dump.sql")
	if err != nil {
		return err, nil
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return err, nil
	}

	db, err := sql.Open("pgx", getConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//TODO: pgxpool を利用して copyNum 以上のコネクションを作成
	//TODO: コネクションを取得する関数をラップして、スキーマを明示的に指定

	ctx := context.Background()
	for i := 0; i < copyNum; i++ {
		copy_content := strings.ReplaceAll(string(content), "hoge", fmt.Sprintf("hoge_%v", i))
		if _, err := db.ExecContext(ctx, copy_content); err != nil {
			return err, nil
		}
	}

	return nil, func() error {
		db, err := sql.Open("pgx", getConnectionString())
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		for i := 0; i < copyNum; i++ {
			fmt.Printf("%v_%v", schema, i)
			if _, err := db.Exec(fmt.Sprintf("DROP SCHEMA %v_%v CASCADE", schema, i)); err != nil {
				return err
			}
		}
		return nil
	}
}

func getConnectionString() string {
	const (
		host     = "localhost"
		port     = 5432
		user     = "admin"
		password = "password"
		dbname   = "db"
	)
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}
