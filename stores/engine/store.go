package engine

import (
	"Project/CarDealearship/models"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/google/uuid"
)

type engineStore struct {
}

// nolint:revive // need not be exported
// New factory function
func New() engineStore {
	return engineStore{}
}

// EngineGetByID is the datastore layer function to get engine by its id
func (s engineStore) EngineGetByID(ctx *gofr.Context, id string) (models.Engine, error) {
	var e models.Engine

	err := ctx.DB().QueryRowContext(ctx, "SELECT * from Engine where id=?;", id).
		Scan(&e.EngineID, &e.Displacement, &e.Cylinders, &e.Range)
	if err != nil {
		return models.Engine{}, err
	}

	return e, nil
}

// EngineCreate is the datastore layer function to create a model of an engine
func (s engineStore) EngineCreate(ctx *gofr.Context, engine *models.Engine) (models.Engine, error) {
	engine.EngineID = uuid.New()

	_, err := ctx.DB().ExecContext(ctx, "INSERT INTO Engine (id,displacement,cylinders,`range`) VALUES(?,?,?,?)",
		engine.EngineID.String(), engine.Displacement, engine.Cylinders, engine.Range)
	if err != nil {
		return models.Engine{}, err
	}

	return *engine, nil
}

// EngineDelete to service layer function to delete the engine from database
func (s engineStore) EngineDelete(ctx *gofr.Context, id string) error {
	_, err := ctx.DB().ExecContext(ctx, "delete from Engine where id=?", id)
	if err != nil {
		return err
	}

	return nil
}

// EngineUpdate is a datastore layer function to update a car record in database
func (s engineStore) EngineUpdate(ctx *gofr.Context, id string, engine *models.Engine) (models.Engine, error) {
	_, err := ctx.DB().ExecContext(ctx, "UPDATE Engine SET displacement=?,cylinders=?,`range`=? WHERE Id=?;",
		engine.Displacement, engine.Cylinders, engine.Range, id)
	if err != nil {
		return models.Engine{}, err
	}

	engine.EngineID, _ = uuid.Parse(id)

	return *engine, nil
}
