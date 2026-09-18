package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("检查 Gemsnote 数据库表结构")

	connStr := "host=127.0.0.1 port=5432 user=gemsnote password=gemsnote dbname=gemsnote sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("打开数据库连接失败:", err)
	}
	defer db.Close()

	// 检查所有表的结构
	tables := []string{"users", "notebooks", "notes", "note_contents", "tags", "files", "attachs", "share_notebooks", "share_notes"}

	for _, table := range tables {
		fmt.Printf("\n=== %s 表结构 ===\n", table)

		rows, err := db.Query(`
			SELECT column_name, data_type, is_nullable, column_default
			FROM information_schema.columns 
			WHERE table_name = $1 
			ORDER BY ordinal_position
		`, table)

		if err != nil {
			fmt.Printf("查询表结构失败: %v\n", err)
			continue
		}

		hasRows := false
		for rows.Next() {
			var columnName, dataType, isNullable, columnDefault sql.NullString
			rows.Scan(&columnName, &dataType, &isNullable, &columnDefault)

			defaultStr := ""
			if columnDefault.Valid {
				defaultStr = " DEFAULT " + columnDefault.String
			}

			nullableStr := ""
			if isNullable.String == "NO" {
				nullableStr = " NOT NULL"
			}

			fmt.Printf("  %-30s %-20s%s%s\n",
				columnName.String,
				dataType.String,
				nullableStr,
				defaultStr)
			hasRows = true
		}
		rows.Close()

		if !hasRows {
			fmt.Println("  表不存在或没有列")
		}

		// 检查主键
		var constraintName, columnName string
		err = db.QueryRow(`
			SELECT tc.constraint_name, kcu.column_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage kcu 
			ON tc.constraint_name = kcu.constraint_name
			WHERE tc.table_name = $1 AND tc.constraint_type = 'PRIMARY KEY'
		`, table).Scan(&constraintName, &columnName)

		if err == nil {
			fmt.Printf("  主键: %s (%s)\n", constraintName, columnName)
		}
	}

	// 检查外键关系
	fmt.Println("\n=== 外键关系 ===")
	rows, err := db.Query(`
		SELECT
			tc.table_name as foreign_table,
			kcu.column_name as foreign_column,
			ccu.table_name as primary_table,
			ccu.column_name as primary_column
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu 
			ON tc.constraint_name = kcu.constraint_name
		JOIN information_schema.constraint_column_usage ccu 
			ON ccu.constraint_name = tc.constraint_name
		WHERE tc.constraint_type = 'FOREIGN KEY'
		ORDER BY tc.table_name, kcu.column_name
	`)

	if err != nil {
		fmt.Printf("查询外键失败: %v\n", err)
	} else {
		defer rows.Close()

		hasRows := false
		for rows.Next() {
			var foreignTable, foreignColumn, primaryTable, primaryColumn string
			rows.Scan(&foreignTable, &foreignColumn, &primaryTable, &primaryColumn)
			fmt.Printf("  %s.%s -> %s.%s\n", foreignTable, foreignColumn, primaryTable, primaryColumn)
			hasRows = true
		}

		if !hasRows {
			fmt.Println("  没有找到外键关系")
		}
	}
}
