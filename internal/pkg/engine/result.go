package engine

// NewResult returns a result
func NewResult(affectedRows int, lastInsertID int, columns []string, rows [][]string) Result {
	return Result{
		affectedRows: affectedRows,
		lastInsertID: lastInsertID,
		columns:      columns,
		rows:         rows,
	}
}

// NewSelectResult returns a select result
func NewSelectResult(affectedRows int, columns []string, rows [][]string) Result {
	return NewResult(affectedRows, -1, columns, rows)
}

// NewInsertResult returns an insert result
func NewInsertResult(affectedRows int, lastInsertID int) Result {
	return NewResult(affectedRows, lastInsertID, nil, nil)
}

// NewUpdateResult returns an update result
func NewUpdateResult(affectedRows int) Result {
	return NewResult(affectedRows, -1, nil, nil)
}

// NewDeleteResult returns a delete result
func NewDeleteResult(affectedRows int) Result {
	return NewResult(affectedRows, -1, nil, nil)
}

// NewEmptyResult returns an empty result
func NewEmptyResult() Result {
	return NewResult(0, -1, nil, nil)
}

// Result - SQL statement result
type Result struct {
	affectedRows int
	lastInsertID int
	rows         [][]string
	columns      []string
}

// AffectedRows returns the affected rows
func (r Result) AffectedRows() int {
	return r.affectedRows
}

// LastInsertId returns the last insert id
func (r Result) LastInsertId() int {
	return r.lastInsertID
}

// Rows returns the rows
func (r Result) Rows() [][]string {
	return r.rows
}

// Columns returns the columns
func (r Result) Columns() []string {
	return r.columns
}
