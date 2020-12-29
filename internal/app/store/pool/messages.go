package pool

import (
	"netcat/internal/app/models"
	"netcat/internal/app/store"
)

// SendMsg ...
func (r *Repo) SendMsg(m *models.Message) error {

	if m == nil {
		return store.ErrInvalidInput
	}

	if m.All {
		for _, val := range r.Pool {
			if err := val.Write(m.Text); err != nil {
				return err
			}
			if err := val.Prefix(); err != nil {
				return err
			}
		}
		return nil
	}

	if err := r.Store.History().AddInto(m); err != nil {
		return err
	}

	m.Text = "\n" + m.Text

	for _, val := range r.Pool {

		if val.Room == m.Room && val.Name != m.Author {
			if err := val.Write(m.Text); err != nil {
				return err
			}

			if err := val.Prefix(); err != nil {
				return err
			}
		}

	}

	return nil
}
