# DocDB

![testing](https://github.com/mattgrunwald/docdb/actions/workflows/test.yml/badge.svg)
[![codecov](https://codecov.io/github/mattgrunwald/docdb/graph/badge.svg?token=VT6LQONNNP)](https://codecov.io/github/mattgrunwald/docdb)

Local document database for native apps. Powered by SQLite.

## Initialization

First, you'll need to open and connect to the database using `New`

```go
db := docdb.New("./app.db", "./db_files")
```

You need to provide the location for the database file (`./app.db`) as well as the location of the directory that will hold the files (`./db_files`). You do not have to create either, just specify their paths.

Docdb stores files in a directory outside of the DB file so that these stored files can be opened and edited without the need for temporary files.

## Using `Doc`s

All operations except for `Delete` return one or more `Doc`s. To open a the file represented by a `Doc`, use `Open`.

```go
file, err := doc.Open()
```

## Querying the DB

`docdb` stores all documents as `Docs`. A `Doc` stores a file's name, it's database ID, when it was added to the DB, and when it was last updated in the DB.

### Insert/Update

```go
// insert a file
doc, err := db.Insert(myFile)
// whoops, we need to use a different file
updatedDoc, err := db.Update(doc.ID, myUpdatedFile)
```

### Find One

Retrieve a `Doc` by its ID:

```go
doc, err := db.FindOne(id)
```

### Find Many

Retrieve many `Doc`s in a specific order

```go
// gets 5 documents with no offset ordered by name in ascending order.
docs, err := db.FindMany(5, 0, col.Name, order.ASC)
```

### Find All

Retrieve all `Doc`s in a specific order

```go
// gets 5 documents with no offset ordered by creation date in descending order.
docs, err := db.FindAll(col.CreatedAt, order.DESC)
```

### Find Like

Retrieve all `Doc`s whose `Name` is similar to the provided search `term`

```go
// returns all documents that have ".log" in their names.
docs, err := db.FindLike(".log")
```

### Delete

Delete a `Doc` using its `ID`

```go
err := db.Delete(doc.ID)
```
