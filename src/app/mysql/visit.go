package mysql

import (
	"context"
	"database/sql"
	"github.com/eldad87/go-boilerplate/src/app"
	"github.com/eldad87/go-boilerplate/src/app/mysql/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
)

func NewVisitService(db *sql.DB) *VisitService {
	return &VisitService{db}
}

type VisitService struct {
	db *sql.DB
}

func (vs *VisitService) Get(c context.Context, id *int) (*app.Visit, error) {
	bVisit, err := models.FindVisit(c, vs.db, *id)
	if err != nil {
		return nil, err
	}

	return sqlBoilerToVisit(bVisit), nil
}

func (vs *VisitService) Set(c context.Context, v *app.Visit) (*app.Visit, error) {
	bVisit := models.Visit{
		ID:        v.ID,
		FirstName: null.StringFrom(v.FirstName),
		LastName:  null.StringFrom(v.LastName),
	}

	err := bVisit.Upsert(c, vs.db, boil.Infer(), boil.Infer())
	if err != nil {
		return nil, err
	}

	return sqlBoilerToVisit(&bVisit), nil
}

func sqlBoilerToVisit(bVisit *models.Visit) *app.Visit {
	return &app.Visit{
		ID:        bVisit.ID,
		FirstName: bVisit.FirstName.String,
		LastName:  bVisit.LastName.String,
		CreatedAt: bVisit.CreatedAt,
		UpdatedAt: bVisit.UpdatedAt,
	}
}
