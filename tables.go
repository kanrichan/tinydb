package tinydb

import "reflect"

// Table the table of the database.
type Table[T any] struct {
	name     string
	database *Database
}

// GetTable get a specific table from the database.
func GetTable[T any](database *Database) *Table[T] {
	var t T
	name := reflect.TypeOf(t).Name()
	return &Table[T]{name, database}
}

// DropTable drop a specific table from the database.
func DropTable[T any](database *Database) error {
	database.Lock()
	defer database.Unlock()
	var t T
	name := reflect.TypeOf(t).Name()
	var data = make(map[string][]interface{})
	err := database.storage.Read(data)
	if err != nil {
		return err
	}
	delete(data, name)
	return database.storage.Write(data)
}

// Insert insert a new document into table.
func (tbl *Table[T]) Insert(items ...T) error {
	tbl.database.Lock()
	defer tbl.database.Unlock()
	var data = make(map[string][]interface{})
	err := tbl.database.storage.Read(&data)
	if err != nil {
		return err
	}
	if _, ok := data[tbl.name]; !ok {
		data[tbl.name] = nil
	}
	for i := range items {
		data[tbl.name] = append(data[tbl.name], items[i])
	}
	return tbl.database.storage.Write(data)
}

// Delete delete all macthing documents in table.
func (tbl *Table[T]) Delete(condition func(T) bool) ([]T, error) {
	tbl.database.Lock()
	defer tbl.database.Unlock()
	var data = make(map[string][]T)
	err := tbl.database.storage.Read(&data)
	if err != nil {
		return nil, err
	}
	var out []T
	if _, ok := data[tbl.name]; !ok {
		return nil, nil
	}
	for i := range data[tbl.name] {
		item := data[tbl.name][i]
		if !condition(item) {
			out = append(out, item)
		}
	}
	var other = make(map[string][]interface{})
	err = tbl.database.storage.Read(&other)
	if err != nil {
		return nil, err
	}
	if _, ok := data[tbl.name]; !ok {
		return nil, nil
	}
	data[tbl.name] = out
	return out, tbl.database.storage.Write(data)
}

// Update update all macthing documents with updater.
func (tbl *Table[T]) Update(updater func(T) T, condition func(T) bool) error {
	tbl.database.Lock()
	defer tbl.database.Unlock()
	var data = make(map[string][]T)
	err := tbl.database.storage.Read(&data)
	if err != nil {
		return err
	}
	var out []T
	if _, ok := data[tbl.name]; !ok {
		return nil
	}
	for i := range data[tbl.name] {
		item := data[tbl.name][i]
		if condition(item) {
			item = updater(item)
		}
		out = append(out, item)
	}
	var other = make(map[string][]interface{})
	err = tbl.database.storage.Read(&other)
	if err != nil {
		return err
	}
	if _, ok := data[tbl.name]; !ok {
		return nil
	}
	data[tbl.name] = out
	return tbl.database.storage.Write(data)
}

// Select select for all documents matching condition.
func (tbl *Table[T]) Select(condition func(T) bool) ([]T, error) {
	var data = make(map[string][]T)
	err := tbl.database.storage.Read(&data)
	if err != nil {
		return nil, err
	}
	var out []T
	if _, ok := data[tbl.name]; !ok {
		return nil, nil
	}
	for i := range data[tbl.name] {
		item := data[tbl.name][i]
		if condition(item) {
			out = append(out, item)
		}
	}
	return out, nil
}
