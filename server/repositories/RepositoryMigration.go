package repositories

type RepositoryMigration interface {
	Migrate() error
}
