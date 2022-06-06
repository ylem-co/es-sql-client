package errors

import "fmt"

const ErrorOperationVersion = "version"
const ErrorOperationSqlQuery = "sql query"

type BadRequest struct{
	Operation string
	Message string
}

func (m *BadRequest) Error() string {
	return fmt.Sprintf("es: %s: bad request %s", m.Operation, m.Message)
}

type Forbidden struct{
	Operation string
	Message string
}

func (m *Forbidden) Error() string {
	return fmt.Sprintf("es: %s: bad request %s", m.Operation, m.Message)
}
