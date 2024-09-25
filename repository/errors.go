package repository

import "errors"

var (
	ErrDatabaseConnectionFailed = errors.New("failed to open a database connection")
	ErrPingFailed               = errors.New("failed to ping database")
	ErrMigrationFailed          = errors.New("failed to run migrations")
	ErrDriverCreationFailed     = errors.New("failed to create a database driver")
	ErrEmbedFailed              = errors.New("failed to create embedded migrations source")
	ErrMigrationInstanceFailed  = errors.New("failed to create a migration instance")
	ErrMigrationRunFailed       = errors.New("migrations failed during run")
	ErrTxRollbackFailed         = errors.New("failed to roll back transaction")
)
