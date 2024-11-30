package engine

func NewResult(affectedRows int, lastInsertId int, columns []string, rows [][]string) Result {
	return Result{
		affectedRows: affectedRows,
		lastInsertID: lastInsertId,
		columns:      columns,
		rows:         rows,
	}
}

func NewSelectResult(affectedRows int, columns []string, rows [][]string) Result {
	return NewResult(affectedRows, -1, columns, rows)
}

func NewInsertResult(affectedRows int, lastInsertId int) Result {
	return NewResult(affectedRows, lastInsertId, nil, nil)
}

func NewUpdateResult(affectedRows int) Result {
	return NewResult(affectedRows, -1, nil, nil)
}

func NewDeleteResult(affectedRows int) Result {
	return NewResult(affectedRows, -1, nil, nil)
}

func NewEmptyResult() Result {
	return NewResult(0, -1, nil, nil)
}

type Result struct {
	affectedRows int
	lastInsertID int
	rows         [][]string
	columns      []string
}

func (r Result) AffectedRows() int {
	return r.affectedRows
}

func (r Result) LastInsertId() int {
	return r.lastInsertID
}

func (r Result) Rows() [][]string {
	return r.rows
}

func (r Result) Columns() []string {
	return r.columns
}
