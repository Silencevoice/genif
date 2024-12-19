package store

type SQLModel interface {
	TableName() string
}

func TableName[T SQLModel]() string {
	var instance T
	return instance.TableName()
}

const SelectAllFieldsFromTable = "SELECT * FROM %s"
const SelectAllFieldsWhereIdEquals = "SELECT * FROM %s WHERE id = ?"
const SelectAllFieldsWhereIdIn = "SELECT * FROM %s WHERE id IN (?)"
const InsertAllFieldsWhereIdEquals = "INSERT INTO %s (%s) VALUES (%s)"
