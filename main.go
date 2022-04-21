package main

import (
	_ "github.com/lib/pq"
)

var schema = `
ALTER TABLE IF EXISTS case_notes DROP COLUMN IF EXISTS outcome_options;
ALTER TABLE IF EXISTS case_notes DROP COLUMN IF EXISTS outcome_rec;

DROP domain IF EXISTS outcome_opts_t_domain;
DROP TYPE IF EXISTS outcome_opts_t;
DROP TABLE IF EXISTS case_notes;

DROP INDEX IF EXISTS ix_case_notes_outcome_options;

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

CREATE INDEX ix_case_notes_outcome_options ON case_notes (outcome_options);
CREATE INDEX ix_case_notes_outcome_options_follow_up ON case_notes (((outcome_options).follow_up));
CREATE INDEX ix_case_notes_outcome_options_hospital_visit_required ON case_notes (((outcome_options).hospital_visit_required));
CREATE INDEX ix_case_notes_outcome_options_prescription_refill_requested ON case_notes (((outcome_options).prescription_refill_requested));
CREATE INDEX ix_case_notes_outcome_options_other ON case_notes (((outcome_options).other));
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
