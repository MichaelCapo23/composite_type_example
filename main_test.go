package main

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func TestQueryData(t *testing.T) {
	db, err := sqlx.Connect("postgres", "user=root password=postgres dbname=michaelcomposite sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)

	stmt := `INSERT INTO case_notes (outcome_options, outcome_rec) VALUES ($1, $2)`
	tx := db.MustBegin()
	for i := 0; i < 50000; i++ {
		stmt := `INSERT INTO case_notes (outcome_options, outcome_rec) VALUES ($1, $2)`
		tx.MustExec(stmt, fmt.Sprintf(`(%t, %t, %t, %t)`, randomBool(), randomBool(), randomBool(), randomBool()), `Get this dog some drugs!`)
	}
	tx.Commit()
	fmt.Println(time.Now().Format(time.RFC3339Nano))
	// fails (compositeStruct & compositeStructBase)
	// var cs []compositeStructBase
	// stmt = `SELECT * FROM case_notes`
	// tx = db.MustBegin()

	// err = tx.Select(&cs, stmt)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	//fails
	// var cs []compositeStruct
	// stmt = `
	// SELECT
	// 	id,
	// 	ROW(
	// 		(outcome_options).follow_up::bool,
	// 		(outcome_options).hospital_visit_required::bool,
	// 		(outcome_options).prescription_refill_requested::bool,
	// 		(outcome_options).other::bool) AS outcome_options,
	// 	outcome_rec
	// FROM
	// 	case_notes`
	// tx = db.MustBegin()

	// err = tx.Select(&cs, stmt)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// works
	var csb int
	stmt = `
	SELECT COUNT(*)
	FROM 
		case_notes
	WHERE
		(outcome_options).follow_up = true
	`
	tx = db.MustBegin()

	err = tx.Get(&csb, stmt)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(time.Now().Format(time.RFC3339Nano))

	// works
	var csb2 []compositeStructBase
	stmt = `
	SELECT 
		id, 
		(outcome_options).*, 
		outcome_rec 
	FROM case_notes`
	tx = db.MustBegin()

	err = tx.Select(&csb2, stmt)
	if err != nil {
		log.Fatalln(err)
	}

	// works
	var csb3 []compositeStructBase
	stmt = `
	SELECT 
		id, 
		(outcome_options).*, 
		outcome_rec 
	FROM 
		case_notes
	WHERE
		(outcome_options).follow_up = true
	`
	tx = db.MustBegin()

	err = tx.Select(&csb3, stmt)
	if err != nil {
		log.Fatalln(err)
	}
}

func randomBool() bool {
	return rand.Intn(2) == 0
}

func TestWhereClauseOneAttribute(t *testing.T) {
	db, err := sqlx.Connect("postgres", "user=root password=postgres dbname=michaelcomposite sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)
	tx := db.MustBegin()
	for i := 0; i < 10000; i++ {
		stmt := `INSERT INTO case_notes (outcome_options, outcome_rec) VALUES ($1, $2)`
		tx.MustExec(stmt, fmt.Sprintf(`(%t, %t, %t, %t)`, randomBool(), randomBool(), randomBool(), randomBool()), `Get this dog some drugs!`)
	}

	var csb []compositeStructBase
	start := time.Now()
	stmt2 := `
			SELECT 
				id, 
				(outcome_options).*, 
				outcome_rec 
			FROM 
				case_notes
			WHERE
				(outcome_options).follow_up = true`

	err = tx.Select(&csb, stmt2)
	end := time.Now()

	totalTime := end.Sub(start)
	fmt.Println(totalTime.Milliseconds())
	if err != nil {
		t.Fatal(err)
	}
}

func TestWhereClauseMultipleAttribute(t *testing.T) {
	db, err := sqlx.Connect("postgres", "user=root password=postgres dbname=michaelcomposite sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	tx := db.MustBegin()
	for i := 0; i < 10000; i++ {
		stmt := `INSERT INTO case_notes (outcome_options, outcome_rec) VALUES ($1, $2)`
		tx.MustExec(stmt, fmt.Sprintf(`(%t, %t, %t, %t)`, randomBool(), randomBool(), randomBool(), randomBool()), `Get this dog some drugs!`)
	}

	var csb []compositeStructBase
	start := time.Now()
	stmt2 := `
			SELECT 
				id, 
				(outcome_options).*, 
				outcome_rec 
			FROM 
				case_notes
			WHERE
				(outcome_options).follow_up = true
			AND
				(outcome_options).hospital_visit_required = true
			AND
				(outcome_options).other = false`

	err = tx.Select(&csb, stmt2)
	end := time.Now()

	totalTime := end.Sub(start)
	fmt.Println(totalTime.Milliseconds())

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(len(csb))
}
