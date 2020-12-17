package utils

const (
	GoInt       = "int"
	GoInt64     = "int64"
	GoTime      = "time.Time"
	GoFloat32   = "float32"
	GoFloat64   = "float64"
	GoByteArray = "[]byte"

	SqlNullInt32   = "sql.NullInt32"
	SqlNullInt64   = "sql.NullInt64"
	SqlNullString  = "sql.NullString"
	SqlNullTime    = "sql.NullTime"
	SqlNullFloat64 = "sql.NullFloat64"

	GormPrimary = ";primary_key"

	JsonKey = "json"
	TargetKey = "target"
	ModeKey = "mode"
	GoFileSuffix = ".go"
)

type Param struct {
	Protocol   string
	UserName   string
	Password   string
	Host       string
	Port       int
	DBName     string
	TableNames []string
	Attachment map[string]interface{}
}

type ColumnInfo struct {
	Name     string
	DataType string
	Nullable bool
	Comment  string
	Key      string
	GoType   string
}

type TableInfo struct {
	Name    string
	Columns []ColumnInfo
}
