package AuditLog

import (
	"database/sql"
	"time"
)

type ConfigType string

type ConfigMapTableHistory struct {
	Id         int           `sql:"id,pk"`
	PipelineId sql.NullInt32 `sql:"pipeline_id"`
	AppId      int           `sql:"app_id"`
	DataType   ConfigType    `sql:"data_type"`
	Data       string        `sql:"data"`
	Deployed   sql.NullBool  `sql:"deployed"`
	DeployedOn *time.Time    `sql:"deployed_on"`
	DeployedBy sql.NullInt32 `sql:"deployed_by"`
	CreatedOn  *time.Time    `sql:"created_on,type:timestamptz"`
	CreatedBy  sql.NullInt32 `sql:"created_by,type:integer"`
	UpdatedOn  *time.Time    `sql:"updated_on,type:timestamptz"`
	UpdatedBy  sql.NullInt32 `sql:"updated_by,type:integer"`
}
