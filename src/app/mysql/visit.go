package mysql

import (
	"context"
	"database/sql"
	"github.com/eldad87/go-boilerplate/src/app"
	"github.com/eldad87/go-boilerplate/src/app/mysql/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"gopkg.in/go-playground/validator.v9"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func NewVisitService(db *sql.DB) *VisitService {
	return &VisitService{db}
}

type VisitService struct {
	db *sql.DB
}

func (vs *VisitService) Get(c context.Context, id *uint) (*app.Visit, error) {
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

	var err error
	if bVisit.ID == 0 {
		err = bVisit.Insert(c, vs.db, boil.Infer())
	} else {
		err = bVisit.Upsert(c, vs.db, boil.Infer(), boil.Infer())
	}

	if err != nil {
		return nil, err
	}

	return sqlBoilerToVisit(&bVisit), nil
}

func (vs *VisitService) Validate(c context.Context, v *app.Visit) error {
	validate = validator.New()
	// TODO: Convert to FieldViolation
	return validate.Struct(v)
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
