module testconfigload/config

require (
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.1.0
	testconfigload/logger v0.0.0
)

replace testconfigload/logger v0.0.0 => ../logger

go 1.12
