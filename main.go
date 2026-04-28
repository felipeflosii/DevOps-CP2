package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Produto struct {
	ID        int64   `json:"id"`
	Nome      string  `json:"nome"`
	Categoria string  `json:"categoria"`
	Preco     float64 `json:"preco"`
}

var db *sql.DB

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		getenv("DB_USER", "root"),
		getenv("DB_PASS", "root"),
		getenv("DB_HOST", "mysql-db"),
		getenv("DB_PORT", "3306"),
		getenv("DB_NAME", "demo"),
	)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/produtos", handleProdutos)
	http.HandleFunc("/produtos/", handleProduto)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	port := getenv("PORT", "8080")
	log.Printf("Servidor rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleProdutos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, nome, categoria, preco FROM produtos")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()
		var produtos []Produto
		for rows.Next() {
			var p Produto
			rows.Scan(&p.ID, &p.Nome, &p.Categoria, &p.Preco)
			produtos = append(produtos, p)
		}
		if produtos == nil {
			produtos = []Produto{}
		}
		json.NewEncoder(w).Encode(produtos)

	case http.MethodPost:
		var p Produto
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		res, err := db.Exec("INSERT INTO produtos (nome, categoria, preco) VALUES (?, ?, ?)", p.Nome, p.Categoria, p.Preco)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		p.ID, _ = res.LastInsertId()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)

	default:
		http.Error(w, "método não permitido", 405)
	}
}

func handleProduto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/produtos/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "id inválido", 400)
		return
	}

	switch r.Method {
	case http.MethodGet:
		var p Produto
		err := db.QueryRow("SELECT id, nome, categoria, preco FROM produtos WHERE id=?", id).
			Scan(&p.ID, &p.Nome, &p.Categoria, &p.Preco)
		if err == sql.ErrNoRows {
			http.Error(w, "não encontrado", 404)
			return
		} else if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(p)

	case http.MethodDelete:
		_, err := db.Exec("DELETE FROM produtos WHERE id=?", id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "método não permitido", 405)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
