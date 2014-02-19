package stores

import (
    "github.com/ledbury/pickleback/elements"
    "github.com/ledbury/pickleback/sets"
    "testing"
    "os"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "crypto/rand"
    "encoding/hex"
    "path/filepath"
)

func tempFile() string {
    randBytes := make([]byte, 16)
    rand.Read(randBytes)
    return filepath.Join(os.TempDir(), "foo"+hex.EncodeToString(randBytes)+".db")
}

func TestInitialize(t *testing.T) {

    filename := tempFile()
    err := Initialize(filename)
    if err != nil {
        t.Error(err)
    }
    defer os.Remove(filename)

    db, _ := sql.Open("sqlite3", filename)
    defer db.Close()
    rows, _ := db.Query("select * from transactions")
    defer rows.Close()
    if cols, _ := rows.Columns(); len(cols) < 1 {
        t.Error("Tables were not created in database")
    }

}

func TestStoreTransaction(t *testing.T) {

    filename := tempFile()
    err := Initialize(filename)
    if err != nil {
        t.Error(err)
    }
    defer os.Remove(filename)

    number, _, err := storeTestTransaction(filename)
    if err != nil {
        t.Error(err)
    }

    db, _ := sql.Open("sqlite3", filename)
    defer db.Close()

    rows, _ := db.Query("select id from transactions")
    defer rows.Close()
    rows.Next()
    var returnedId string
    err = rows.Scan(&returnedId)

    if err != nil {
        t.Error(err)
    }

    if returnedId != number {
        t.Error("Transaction was not stored successfully. Expected", number, "got", returnedId)
    }

}

func TestFindTransaction(t *testing.T) {

    filename := tempFile()
    err := Initialize(filename)
    if err != nil {
        t.Error(err)
    }
    defer os.Remove(filename)

    number, elems, err := storeTestTransaction(filename)
    if err != nil {
        t.Error(err)
    }

    trans, err := FindTransaction(filename, number)
    if err != nil {
        t.Error(err)
    }

    for i := range trans.Elements {
        if elems[i].Id != trans.Elements[i].Id {
            t.Error("Retrieved transaction did not match stored transaction.")
        }
    }

}

func TestElementsForTransaction(t *testing.T) {

    filename := tempFile()
    err := Initialize(filename)
    if err != nil {
        t.Error(err)
    }
    defer os.Remove(filename)

    number, elems, err := storeTestTransaction(filename)
    if err != nil {
        t.Error(err)
    }

    db, err := sql.Open("sqlite3", filename)
    if err != nil { return }
    defer db.Close()
    retrievedElems, err := elementsForTransaction(filename, db, number)
    if err != nil {
        t.Error(err)
    }

    for i := range retrievedElems {
        if elems[i].Id != retrievedElems[i].Id {
            t.Error("Retrieved elements did not match stored elements")
        }
    }

}

func storeTestTransaction(filename string) (string, []*elements.Element, error) {
    number := "R1234"
    elems := []*elements.Element{&elements.Element{Id: int64(123), Name: "Abc"}, &elements.Element{Id: int64(456), Name: "Def"}}
    trans := sets.Transaction{Id: number, Set: sets.Set{Elements: elems}}
    _, err := StoreTransaction(filename, &trans)
    return number, elems, err
}
