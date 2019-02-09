package main

/*
#ctype Stmt *
*/
type stmtHandle uintptr

/*
#cmethod Open
#cmethod Close
*/
type dbIf struct {
	handle stmtHandle
	dbName string
}

/*
Enum type for operand
#ctype operKind
enum operKind: int32_t {
  Get = 0,
  Put = 1,
  Delete = 2
};
*/
type OperKind int32

const (
	Get    = OperKind(0)
	Put    = OperKind(1)
	Delete = OperKind(2)
)

/*
 */
type dbOper struct {
	kind  OperKind
	key   string
	value []byte
}

/*
#cmethod Do
*/
type dbBatch struct {
	operations []dbOper
}
