package tinydb

type Table[T any] struct {
	name     string
	database *database
}

func (tbl *Table[T]) Insert(items ...T) error {
	tbl.database.Lock()
	defer tbl.database.Lock()
	var data = make(map[string][]interface{})
	err := tbl.database.storage.Read(data)
	if err != nil {
		return err
	}
	for i := range items {
		data[tbl.name] = append(data[tbl.name], items[i])
	}
	return tbl.database.storage.Write(data)
}

func (tbl *Table[T]) Select(condition func(T) bool) ([]T, error) {
	var data = make(map[string][]T)
	err := tbl.database.storage.Read(&data)
	if err != nil {
		return nil, err
	}
	var out []T
	for i := range data[tbl.name] {
		item := data[tbl.name][i]
		if condition(item) {
			out = append(out, item)
		}
	}
	return out, nil
}
