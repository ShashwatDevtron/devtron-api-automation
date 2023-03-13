package AuditLog

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type PipelineStrategyTableHistory struct {
	Id                  int           `sql:"id,pk"`
	PipelineId          int           `sql:"pipeline_id, notnull"`
	Strategy            string        `sql:"strategy,notnull"`
	Config              string        `sql:"config"`
	Default             bool          `sql:"default,notnull"`
	Deployed            sql.NullBool  `sql:"deployed"`
	DeployedOn          pq.NullTime   `sql:"deployed_on"`
	DeployedBy          sql.NullInt32 `sql:"deployed_by"`
	CreatedOn           *time.Time    `sql:"created_on,type:timestamptz"`
	CreatedBy           sql.NullInt32 `sql:"created_by,type:integer"`
	UpdatedOn           *time.Time    `sql:"updated_on,type:timestamptz"`
	UpdatedBy           sql.NullInt32 `sql:"updated_by,type:integer"`
	PipelineTriggerType string        `sql:"pipeline_trigger_type"`
}
