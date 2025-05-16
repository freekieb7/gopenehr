package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/acontrolfreak/openehr/aql"
	"github.com/acontrolfreak/openehr/aql/gen"
	"github.com/antlr4-go/antlr/v4"

	_ "github.com/lib/pq" // add this
)

type AqlResultColumn struct {
	Name string `json:"name"`
	Path string `json:"path,omitempty"`
}

type AqlResult struct {
	Query   string            `json:"q"`
	Columns []AqlResultColumn `json:"columns"`
	Rows    []any             `json:"rows"`
}

type QueryAqlBody struct {
	Query      string         `json:"q"`
	Parameters map[string]any `json:"query_parameters"`
}

func main() {
	connStr := os.Getenv("DB_CONN_STR")
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/aql/push/ehr", func(w http.ResponseWriter, r *http.Request) {
		var jsn map[string]interface{}

		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&jsn); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jsnBytes, _ := json.Marshal(jsn)

		if _, err := db.Exec("INSERT INTO ehr(data) VALUES ($1) ", jsnBytes); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/aql/push/ehr/{ehr_uid}/ehr_status", func(w http.ResponseWriter, r *http.Request) {
		ehrUid := r.PathValue("ehr_uid")

		var ehrId int
		if err := db.QueryRow("SELECT id FROM ehr WHERE uid = $1", ehrUid).Scan(&ehrId); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var data map[string]interface{}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		dataBlob, _ := json.Marshal(data)

		if _, err := db.Exec("INSERT INTO ehr_status(ehr_id, data) VALUES ($1, $2) ", ehrId, dataBlob); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/aql/push/ehr/{ehr_uid}/ehr_access", func(w http.ResponseWriter, r *http.Request) {
		ehrUid := r.PathValue("ehr_uid")

		var ehrId int
		if err := db.QueryRow("SELECT id FROM ehr WHERE uid = $1", ehrUid).Scan(&ehrId); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var data map[string]interface{}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dataBlob, _ := json.Marshal(data)
		if _, err := db.Exec("INSERT INTO ehr_access(ehr_id, data) VALUES ($1, $2) ", ehrId, dataBlob); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/aql/push/ehr/{ehr_uid}/composition", func(w http.ResponseWriter, r *http.Request) {
		ehrUid := r.PathValue("ehr_uid")

		var ehrId int
		if err := db.QueryRow("SELECT id FROM ehr WHERE uid = $1", ehrUid).Scan(&ehrId); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var data map[string]interface{}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dataBlob, _ := json.Marshal(data)
		if _, err := db.Exec("INSERT INTO composition(ehr_id, data) VALUES ($1, $2) ", ehrId, dataBlob); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/ehr", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// TODO
		}

		if r.Method == "POST" {
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/query/aql", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// TODO
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var body QueryAqlBody
		d := json.NewDecoder(r.Body)

		if err := d.Decode(&body); err != nil {
			log.Fatal(err)
		}

		is := antlr.NewInputStream(body.Query)

		lexer := gen.NewAqlLexer(is)
		stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

		// Create the Parser
		parser := gen.NewAqlParser(stream)
		pel := aql.NewErrorListener()
		parser.AddErrorListener(pel)

		// Finally parse the expression
		tree := parser.SelectQuery()

		if pel.Count() > 0 {
			errMsg := strings.Join(pel.Errors, ";")

			if _, err := w.Write([]byte(errMsg)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		listener := aql.NewAqlParserListener()
		antlr.ParseTreeWalkerDefault.Walk(listener, tree)

		builder := aql.NewQueryBuilder(listener, body.Parameters)
		query, err := builder.Build()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		var res AqlResult

		rows, err := db.Query(query)
		log.Println(query)

		if err != nil {
			// handle err
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Query = body.Query
		res.Rows = make([]any, 0)

		// TODO
		for i, sel := range listener.Selects() {
			var col AqlResultColumn

			if sel.IDENTIFIER() != nil {
				col.Name = sel.IDENTIFIER().GetText()
			} else {
				col.Name = fmt.Sprintf("#%d", i)
			}

			col.Path = sel.ColumnExpr().GetText()

			res.Columns = append(res.Columns, col)
		}

		colTypes, _ := rows.ColumnTypes()
		ptrs := make([]any, len(colTypes))
		vals := make([]reflect.Value, len(colTypes))
		for i, colType := range colTypes {
			var v reflect.Value

			log.Println(colType.DatabaseTypeName())

			switch colType.DatabaseTypeName() {
			case "INT4":
				{
					v = reflect.New(reflect.TypeOf(new(int8)))
				}

			case "INT8":
				{
					v = reflect.New(reflect.TypeOf(new(int16)))
				}
			case "NUMERIC":
				{
					v = reflect.New(reflect.TypeOf(new(float64)))
				}
			case "JSONB":
				{
					v = reflect.New(reflect.TypeOf(new(json.RawMessage)))
				}
			default:
				{
					v = reflect.New(reflect.TypeOf(new(string)))
				}
			}

			//vals[i] = reflect.New(colType.ScanType())
			vals[i] = v
			ptrs[i] = vals[i].Interface()
		}

		result := make([]any, len(colTypes))

		for rows.Next() {
			if err := rows.Scan(ptrs...); err != nil {
				log.Fatal(err)
				return
			}

			for i, v := range vals {
				result[i] = v.Elem().Interface()
			}

			res.Rows = append(res.Rows, result)
		}

		//
		//for rows.Next() {
		//	err = rows.Scan(pointers...)
		//
		//	for i, v := range vals {
		//		result[names[i]] = v.Elem().Interface()
		//	}
		//
		//	res.Rows = append(res.Rows, pointers)
		//	break
		//}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//w.WriteHeader(http.StatusOK)
	})

	// And we serve HTTP until the world ends.
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
