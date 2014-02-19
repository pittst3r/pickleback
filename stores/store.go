package stores

import (
    "github.com/ledbury/pickleback/elements"
    "github.com/ledbury/pickleback/sets"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "os"
    // "errors"
)

func Initialize(dbFilename string) error {

    os.Remove(dbFilename)

    db, err := sql.Open("sqlite3", dbFilename)
    defer db.Close()

    trans := `
    create table transactions (
        id              text not null primary key
    );
    `
    elems := `
    create table elements (
        id              integer not null primary key autoincrement,
        uid             integer not null,
        name            text,
        transaction_id  text not null
    );
    `

    _, err = db.Exec(trans)
    _, err = db.Exec(elems)

    return err

}

// Returns id of transaction
func StoreTransaction(dbFilename string, transaction *sets.Transaction) error {

    db, err := sql.Open("sqlite3", dbFilename)
    if err != nil { return err }
    defer db.Close()
    tx, _ := db.Begin()

    // Insert transaction
    _, err = tx.Exec("insert into transactions(id) values(?)", transaction.Id)
    if err != nil { return err }

    // Insert elements
    stmt, err := tx.Prepare("insert into elements(uid, name, transaction_id) values(?, ?, ?)")
    if err != nil { return err }
    defer stmt.Close()
    for _, e := range transaction.Elements {
        _, err = stmt.Exec(e.Id, e.Name, transaction.Id)
        if err != nil { return err }
    }

    tx.Commit()

    return err

}

func FindTransaction(dbFilename string, id string) (trans *sets.Transaction, err error) {

    trans = &sets.Transaction{Id: id}

    db, err := sql.Open("sqlite3", dbFilename)
    if err != nil { return }
    defer db.Close()

    elems, err := elementsForTransaction(dbFilename, db, id)
    if err != nil { return }
    trans.Elements = elems

    return
}

func elementsForTransaction(dbFilename string, db *sql.DB, id string) (elems []*elements.Element, err error) {

    rows, err := db.Query("select uid, name from elements where transaction_id = ?", id)
    if err != nil { return }
    defer rows.Close()
    for rows.Next() {
        var eId int
        var eName string
        err = rows.Scan(&eId, &eName)
        elems = append(elems, &elements.Element{Id: int64(eId), Name: eName})
    }
    
    return elems, err

}
