# tinydb

### Introduction

[TinyDB](https://github.com/msiemens/tinydb) is a lightweight document oriented database optimized for your happiness :) The target are small apps that would be blown away by a SQL-DB or an external database server.

But in this project, it's written in pure Golang and has no external dependencies. 

### Example Code

Import
```Go
import tiny "github.com/Yiwen-Chan/tinydb"
```

New a storage
```Go
storage, err := tiny.JSONStorage("test.json")
```

Open a database used by storage
```Go
database, err := tiny.TinyDB(storage)
defer database.Close()
```

Get a table and Insert or Delete or Update or Select
```Go
table := tiny.GetTable[student](database)
table.Insert(student{001, "test"})
table.Update(func(s student) student { s.ID = 002; return s }, func(s student) bool { return true })
table.Select(func(s student) bool { return true })
table.Delete(func(s student) bool { return true })
```
