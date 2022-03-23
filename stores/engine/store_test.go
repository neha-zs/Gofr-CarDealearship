package engine

import (
	"Project/CarDealearship/models"
	"context"
	"database/sql"
	"fmt"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

// TestEngineGetByID test the EngineGetByID functionality of the datastore layer
func TestEngineGetByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.TODO()

	dbcheck := New()

	defer db.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	engine := models.Engine{EngineID: id, Displacement: 1800, Cylinders: 7, Range: 0}

	rows := sqlmock.NewRows([]string{"id", "displacement", "cylinders", "range"}).
		AddRow(id.String(), 1800, 7, 0)
	mock.ExpectQuery("SELECT * from Engine where id=?;").WithArgs(id.String()).WillReturnRows(rows)
	mock.ExpectQuery("SELECT * from Engine where id=?;").WithArgs(uuid.Nil).WillReturnError(errors.EntityNotFound{})

	cases := []struct {
		desc   string
		input  uuid.UUID
		output models.Engine
		err    error
	}{
		{"success", engine.EngineID, engine, nil},
		{"failure", uuid.Nil, models.Engine{}, errors.EntityNotFound{}},
	}
	for i, tc := range cases {
		resp, err := dbcheck.EngineGetByID(ctx, tc.input.String())

		if resp != tc.output {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.output)
		}

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

// TestEngineCreate test the EngineCreate functionality of the datastore layer
func TestEngineCreate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.TODO()

	dbcheck := New()

	defer db.Close()

	id := uuid.New()
	engine := models.Engine{EngineID: id, Displacement: 1600, Cylinders: 4, Range: 0}

	mock.ExpectExec("INSERT INTO Engine (id,displacement,cylinders,`range`) VALUES(?,?,?,?)").
		WithArgs(sqlmock.AnyArg(), engine.Displacement, engine.Cylinders, engine.Range).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO Engine (id,displacement,cylinders,`range`) VALUES(?,?,?,?)").
		WithArgs(sqlmock.AnyArg(), engine.Displacement, engine.Cylinders, engine.Range).
		WillReturnError(errors.Error("Entry Failed"))

	cases := []struct {
		desc string
		id   uuid.UUID
		err  error
	}{
		{"success", engine.EngineID, nil},
		{"failure", uuid.Nil, errors.Error("Entry Failed")},
	}

	for i, tc := range cases {
		_, err := dbcheck.EngineCreate(ctx, &engine)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

// TestEngineUpdate test the EngineUpdate functionality of the datastore layer
func TestEngineUpdate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.TODO()

	dbcheck := New()

	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(db)

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	engine := models.Engine{EngineID: id, Displacement: 1800, Cylinders: 8, Range: 1}

	mock.ExpectExec("UPDATE Engine SET displacement=?,cylinders=?,`range`=? WHERE Id=?;").
		WithArgs(engine.Displacement, engine.Cylinders, engine.Range, engine.EngineID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE Engine SET displacement=?,cylinders=?,`range`=? WHERE Id=?;").
		WithArgs(engine.Displacement, engine.Cylinders, engine.Range, engine.EngineID).
		WillReturnError(errors.EntityNotFound{})

	cases := []struct {
		desc  string
		input models.Engine
		err   error
	}{
		{"success", engine, nil},
		{"failure", engine, errors.EntityNotFound{}},
	}

	for i, tc := range cases {
		_, err := dbcheck.EngineUpdate(ctx, tc.input.EngineID.String(), &tc.input)
		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

// TestEngineDelete tests for EngineDelete  function
func TestEngineDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.TODO()

	dbcheck := New()

	defer db.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	mock.ExpectExec("delete from Engine where id=?").WithArgs(id.String()).WillReturnResult(sqlmock.NewResult(
		1, 1))
	mock.ExpectExec("delete  from Engine where id=?").WithArgs(uuid.Nil).WillReturnError(errors.EntityNotFound{})

	cases := []struct {
		desc string
		id   uuid.UUID
		err  error
	}{
		{"Delete success ", id, nil},
		{"Delete failed", uuid.Nil, errors.EntityNotFound{}},
	}

	for i, tc := range cases {
		err := dbcheck.EngineDelete(ctx, tc.id.String())
		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}
