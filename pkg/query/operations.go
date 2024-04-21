package query

import (
	"fmt"

	"strings"
)

func AddSelect(tableName string, column string) {
	// str := fmt.Sprintf("%s.%s as %s_%s", tableName, column, formatName(tableName), column)
	// containes := false
	// for i := range s.Rows {
	// 	if strings.EqualFold(s.Rows[i], str) {
	// 		containes = true
	// 	}
	// }
	// if !containes {
	// 	s.Rows = append(s.Rows, str)
	// }
}

func AddJoin(mainTableName string, mainTableCol string, connectingTableName string, connectingTableCol string, builder *strings.Builder) {
	builder.WriteString(" inner join ")
	builder.WriteString(connectingTableName)
	builder.WriteString(" on ")
	builder.WriteString(mainTableName + "." + mainTableCol)
	builder.WriteString(" = ")
	builder.WriteString(connectingTableName + "." + connectingTableCol)
}

type jc struct {
	TableName           string
	ColumnName          string
	ReferenceTableName  string
	ReferenceColumnName string
}

func (j *jc) isSameCol(col jc) bool {
	return j.TableName == col.TableName &&
		j.ColumnName == col.ColumnName &&
		j.ReferenceTableName == col.ReferenceTableName &&
		j.ReferenceColumnName == col.ReferenceColumnName
}

type Join struct {
	Type   string
	Column jc
}

func (q *Query) RemoveJoin(join Join) {
	if len(q.Joins) == 0 {
		return
	}
	jc := jc{
		TableName:           join.Column.TableName,
		ColumnName:          join.Column.ColumnName,
		ReferenceTableName:  join.Column.ReferenceTableName,
		ReferenceColumnName: join.Column.ReferenceColumnName,
	}

	for i, v := range q.Joins {
		if v.Column.isSameCol(jc) {
			q.Joins = append(q.Joins[:i], q.Joins[i+1:]...)
			return
		}
	}
}

func (q *Query) AddJoin(
	Type string,
	tableName string,
	columnName string,
	referenceTableName string,
	referenceColumnName string,
) {
	jc := jc{
		TableName:           tableName,
		ColumnName:          columnName,
		ReferenceTableName:  referenceTableName,
		ReferenceColumnName: referenceColumnName,
	}
	for _, v := range q.Joins {
		isEqual := strings.Compare(v.Column.TableName, tableName)
		fmt.Println(isEqual)
		if isEqual == 0 {
			return
		}
	}
	q.Joins = append(q.Joins, Join{
		Type:   Type,
		Column: jc,
	})

}

type Query struct {
	Rows    []string
	Headers []string
	Joins   []Join
}

func NewQuery() *Query {
	return &Query{
		Rows:    make([]string, 0),
		Headers: make([]string, 0),
		Joins:   make([]Join, 0),
	}
}

func formatName(name string) string {
	splitted := strings.Split(name, "_")
	full_str := ""

	for i := range splitted {
		full_str += splitted[i]
	}
	return full_str
}

func (s *Query) AddSelect(tableName string, column string) {
	str := fmt.Sprintf("%s.%s as %s_%s", tableName, column, formatName(tableName), column)
	containes := false
	for i := range s.Rows {
		if strings.EqualFold(s.Rows[i], str) {
			containes = true
		}
	}
	if !containes {
		s.Rows = append(s.Rows, str)
	}
}

func (s *Query) AddHeader(header string) error {
	containes := false
	for i := range s.Headers {
		if strings.Contains(s.Headers[i], header) {
			containes = true
		}
	}

	if !containes {
		s.Headers = append(s.Headers, header)
	}
	return nil
}
func (s *Query) GetRows() (*[]string, error) {
	return &s.Rows, nil
}
