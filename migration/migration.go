package migration

import "time"

type Migrator interface {
	Connect() error
	Close() error
	PreMigrate() error
	MigrateTable(tableName string) (MigrationResult, error)
	PostMigrate() error
	Validate() (ValidationResult, error)
}

type MigrationResult struct {
	TableName      string
	TotalRecords   int
	SuccessRecords int
	FailedRecords  int
	Errors         []string
	Duration       time.Duration
}

type MigrationProgress struct {
	CurrentTable    string
	Progress        float64
	TotalRecords    int
	MigratedRecords int
	CurrentRecord   string
}

type ValidationResult struct {
	Table           string
	SourceCount     int
	TargetCount     int
	RecordsMatch    bool
	MissingInTarget []string
	ExtraInTarget   []string
	DataIntegrity   []string
	Passed          bool
}

type MigrationConfig struct {
	BatchSize     int
	StopOnError   bool
	ValidateAfter bool
	Verbose       bool
}

type MigratorBase struct {
	Config MigrationConfig
}

func NewMigratorBase(config MigrationConfig) *MigratorBase {
	if config.BatchSize <= 0 {
		config.BatchSize = 1000
	}
	return &MigratorBase{Config: config}
}
