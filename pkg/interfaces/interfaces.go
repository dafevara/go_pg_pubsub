package interfaces

import (
	"time"

	"github.com/google/uuid"
)

type Taskable interface {
	Next() (int64, error)
	Perform() error
}

type BaseTaskable struct {
	Uuid       uuid.UUID `pg:"DEFAULT:gen_random_uuid()"`
	QueueName  string
	TriesLeft  int32 `pg:"default:5"`
	Error      string
	Processing bool `pg:"default:false"`
	NextTryAt  time.Time
	UpdatedAt  time.Time
}

func (base *BaseTaskable) New() error {
	return nil
}

func (base *BaseTaskable) Next() (int64, error) {
	return 1, nil
}
