package memory

import (
	"errors"

	"github.com/cdevoogd/dashboard/backend/dash"
)

var errNotFound = errors.New("object not found")

type Store struct {
	apps map[string]*dash.ApplicationRecord
}

func NewStore() *Store {
	return &Store{
		apps: make(map[string]*dash.ApplicationRecord),
	}
}

func (s *Store) IsNotFoundError(err error) bool {
	return errors.Is(err, errNotFound)
}

func (s *Store) AddApplication(app *dash.ApplicationRecord) error {
	s.apps[app.ID] = app
	return nil
}

func (s *Store) GetApplication(id string) (*dash.ApplicationRecord, error) {
	app, ok := s.apps[id]
	if !ok {
		return nil, errNotFound
	}
	return app, nil
}

func (s *Store) GetAllApplications() ([]*dash.ApplicationRecord, error) {
	var apps []*dash.ApplicationRecord
	for _, app := range s.apps {
		apps = append(apps, app)
	}
	return apps, nil
}

func (s *Store) UpdateApplication(app *dash.ApplicationRecord) error {
	_, ok := s.apps[app.ID]
	if !ok {
		return errNotFound
	}

	s.apps[app.ID] = app
	return nil
}

func (s *Store) DeleteApplication(id string) error {
	_, ok := s.apps[id]
	if !ok {
		return errNotFound
	}

	delete(s.apps, id)
	return nil
}
