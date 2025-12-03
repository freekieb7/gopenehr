package aql

// type QueryHandler struct {
// 	logger *slog.Logger
// 	db     *storage.Database
// }

// func NewQueryHandler(logger *slog.Logger, db *storage.Database) QueryHandler {
// 	return QueryHandler{logger: logger, db: db}
// }

// type ExecuteQueryRequest struct {
// 	AQL        string         `json:"aql"`
// 	Parameters map[string]any `json:"parameters,omitempty"`
// }

// func (h *QueryHandler) ExecuteQuery(w http.ResponseWriter, r *http.Request) error {
// 	ctx := r.Context()

// 	// Parse request body
// 	var req ExecuteQueryRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		return err
// 	}

// 	queryContext, err := aql.QueryContext(req.AQL)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return nil
// 	}

// 	knownPreparedTables, err := h.db.GetAllPreparedTables(ctx) // Make sure prepared tables are loaded
// 	if err != nil {
// 		h.logger.Error("get prepared tables error", "error", err)
// 		http.Error(w, "internal server error", http.StatusInternalServerError)
// 		return nil
// 	}

// 	preparedTables := make([]aql.PreparedTable, len(knownPreparedTables))
// 	for i, query := range knownPreparedTables {
// 		ctx, err := aql.QueryContext(query.AQL)
// 		if err != nil {
// 			h.logger.Error("parse prepared table error", "error", err)
// 			http.Error(w, "internal server error", http.StatusInternalServerError)
// 			return nil
// 		}

// 		preparedTables[i] = aql.PreparedTable{
// 			Name:   query.Name,
// 			Source: query.TableName,
// 			Ctx:    ctx.SelectQuery(),
// 		}
// 	}

// 	q, cols, err := aql.BuildQuery(queryContext, req.Parameters, preparedTables)
// 	if err != nil {
// 		h.logger.Error("build query error", "error", err)
// 		if buildError, ok := err.(aql.BuildError); ok {
// 			http.Error(w, buildError.Message, http.StatusBadRequest)
// 			return nil
// 		}

// 		h.logger.Error("internal error", "error", err)
// 		http.Error(w, "internal server error", http.StatusInternalServerError)
// 		return nil
// 	}

// 	// Execute query
// 	rows, err := h.db.Query(ctx, q)
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()

// 	colNames := ""
// 	colTypes := ""
// 	for i, col := range cols {
// 		if i > 0 {
// 			colNames += ","
// 			colTypes += ","
// 		}
// 		colName := col.Name
// 		if colName == "" {
// 			colName = fmt.Sprintf("f%d", i)
// 		}
// 		colNames += colName
// 		colTypes += col.Type.Name()
// 	}
// 	w.Header().Set("X-Column-Names", colNames)
// 	w.Header().Set("X-Column-Types", colTypes)
// 	w.Header().Set("Content-Type", "application/json")

// 	w.Write([]byte("["))

// 	first := true
// 	for rows.Next() {
// 		var jsonData []byte
// 		if err := rows.Scan(&jsonData); err != nil {
// 			h.logger.Error("scan error", "error", err)
// 			continue
// 		}

// 		if !first {
// 			w.Write([]byte(","))
// 		}
// 		w.Write(jsonData)
// 		first = false

// 		// Flush each row so client receives data progressively
// 		if f, ok := w.(http.Flusher); ok {
// 			f.Flush()
// 		}
// 	}

// 	w.Write([]byte("]"))
// 	return nil
// }

// type CreatePreparedTableRequest struct {
// 	Name string `json:"name"`
// 	AQL  string `json:"aql"`
// }

// func (h *QueryHandler) CreatePreparedTable(w http.ResponseWriter, r *http.Request) error {
// 	ctx := r.Context()

// 	// Parse request body
// 	var req CreatePreparedTableRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		return err
// 	}

// 	if req.Name == "" {
// 		http.Error(w, "name is required", http.StatusBadRequest)
// 		return nil
// 	}
// 	if req.AQL == "" {
// 		http.Error(w, "aql is required", http.StatusBadRequest)
// 		return nil
// 	}

// 	// check if name can be a identifier
// 	for _, r := range req.Name {
// 		if !(r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
// 			http.Error(w, "name can only contain letters, numbers and underscores", http.StatusBadRequest)
// 			return nil
// 		}
// 	}

// 	// Parse AQL to validate
// 	listener := aql.NewTreeShapeListener()
// 	errorListener := aql.NewErrorListener()

// 	input := antlr.NewInputStream(req.AQL)
// 	lexer := gen.NewAQLLexer(input)
// 	stream := antlr.NewCommonTokenStream(lexer, 0)

// 	p := gen.NewAQLParser(stream)
// 	p.AddErrorListener(errorListener)

// 	antlr.ParseTreeWalkerDefault.Walk(listener, p.Query())

// 	if len(errorListener.Errors) > 0 {
// 		http.Error(w, errors.Join(errorListener.Errors...).Error(), http.StatusBadRequest)
// 		return nil
// 	}

// 	sqlQuery, cols, err := aql.BuildSelectQuery(listener.Query.SelectQuery(), make(aql.Parameters), make([]aql.PreparedTable, 0))
// 	if err != nil {
// 		h.logger.Error("build query error", "error", err)
// 		if buildError, ok := err.(aql.BuildError); ok {
// 			http.Error(w, buildError.Message, http.StatusBadRequest)
// 			return nil
// 		}

// 		h.logger.Error("internal error", "error", err)
// 		http.Error(w, "internal server error", http.StatusInternalServerError)
// 		return nil
// 	}

// 	// For prepared queries, all the columns must be named
// 	for i, col := range cols {
// 		if col.Name == "" {
// 			http.Error(w, fmt.Sprintf("all selected columns must be named for prepared queries; column %d is not named", i+1), http.StatusBadRequest)
// 			return nil
// 		}
// 	}

// 	// Validate if executable
// 	if _, err := h.db.Exec(ctx, "EXPLAIN "+sqlQuery); err != nil {
// 		http.Error(w, "invalid query", http.StatusBadRequest)
// 		return nil
// 	}

// 	// Store prepared table
// 	if _, err := h.db.CreatePreparedTable(ctx, storage.CreatePreparedTableParams{
// 		Name: req.Name,
// 		AQL:  req.AQL,
// 		SQL:  sqlQuery,
// 	}); err != nil {
// 		h.logger.Error("create prepared table error", "error", err)
// 		http.Error(w, "internal server error", http.StatusInternalServerError)
// 		return nil
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	return nil
// }

// func (h *QueryHandler) SyncPreparedTable(w http.ResponseWriter, r *http.Request) error {
// 	ctx := r.Context()

// 	name := r.URL.Query().Get("name")

// 	preparedTable, err := h.db.GetPreparedTableByName(ctx, name)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrNoRows) {
// 			http.Error(w, "prepared table not found", http.StatusNotFound)
// 			return nil
// 		}
// 		h.logger.Error("get prepared table error", "error", err)
// 		http.Error(w, "internal server error", http.StatusInternalServerError)
// 		return nil
// 	}

// 	if err := h.db.SyncPreparedTable(ctx, preparedTable.ID); err != nil {
// 		h.logger.Error("sync prepared table error", "error", err)
// 		http.Error(w, "internal server error", http.StatusInternalServerError)
// 		return nil
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// 	return nil
// }
