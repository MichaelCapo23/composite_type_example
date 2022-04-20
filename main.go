package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
ALTER TABLE IF EXISTS case_notes DROP COLUMN IF EXISTS outcome_options;
ALTER TABLE IF EXISTS case_notes DROP COLUMN IF EXISTS outcome_rec;

DROP domain IF EXISTS outcome_opts_t_domain;
DROP TYPE IF EXISTS outcome_opts_t;
DROP TABLE IF EXISTS case_notes;

CREATE TYPE outcome_opts_t AS (
    follow_up boolean,
    hospital_visit_required boolean,
    prescription_refill_requested boolean,
    other boolean 
);

CREATE TABLE case_notes (
    id serial primary key, 
    outcome_options outcome_opts_t,
    outcome_rec text
);

create domain outcome_opts_t_domain as outcome_opts_t 
check (
  (value).follow_up is not null AND
  (value).hospital_visit_required is not null AND
  (value).prescription_refill_requested is not null AND
  (value).other is not null
);

ALTER TABLE case_notes ALTER COLUMN outcome_options TYPE outcome_opts_t_domain;
`

type compositeStructBase struct {
	ID                          int    `db:"id"`
	FollowUp                    bool   `db:"follow_up"`
	HospitalVisitRequired       bool   `db:"hospital_visit_required"`
	PrescriptionRefillRequested bool   `db:"prescription_refill_requested"`
	Other                       bool   `db:"other"`
	OutcomeRec                  string `db:"outcome_rec"`
}

type compositeStructArray struct {
	ID             int     `db:"id"`
	OutcomeOptions []uint8 `db:"outcome_options"`
	OutcomeRec     string  `db:"outcome_rec"`
}

type compositeStruct struct {
	ID             int            `db:"id"`
	OutcomeOptions outcomeOptions `db:"outcome_options"`
	OutcomeRec     string         `db:"outcome_rec"`
}

type outcomeOptions struct {
	FollowUp                    bool `db:"follow_up"`
	HospitalVisitRequired       bool `db:"hospital_visit_required"`
	PrescriptionRefillRequested bool `db:"prescription_refill_requested"`
	Other                       bool `db:"other"`
}

func main() {
	db, err := sqlx.Connect("postgres", "user=root password=postgres dbname=michaelcomposite sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schema)

	stmt := `INSERT INTO case_notes (outcome_options, outcome_rec) VALUES ($1, $2)`
	tx := db.MustBegin()
	tx.MustExec(stmt, `(true,false,true,false)`, `Get this dog some drugs!`)
	tx.MustExec(stmt, `(false,false,false,false)`, `Kitty appears completely healthy`)
	tx.MustExec(stmt, `(false,false,false,false)`, `Gerbil had a temporary virus`)
	tx.MustExec(stmt, `(false,true,false,false)`, `Surgery should be scheduled asap`)
	tx.Commit()

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
	var csb []compositeStructBase
	stmt = `
	SELECT
		id,
		(outcome_options).follow_up,
		(outcome_options).hospital_visit_required,
		(outcome_options).prescription_refill_requested,
		(outcome_options).other,
		outcome_rec
	FROM case_notes`
	tx = db.MustBegin()

	err = tx.Select(&csb, stmt)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(csb)

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

	fmt.Println(csb2)
}
