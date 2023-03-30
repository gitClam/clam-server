package config

type config struct {
	System  system  `mapstructure:"system" yaml:"system"`
	Zap     zap     `mapstructure:"zap" yaml:"zap"`
	Jwt     jwt     `mapstructure:"jwt" yaml:"jwt"`
	Decoder decoder `mapstructure:"decoder" yaml:"decoder"`
}

type system struct {
	Host string `mapstructure:"host" yaml:"host"` // 端口号
}

type zap struct {
	Level         string   `mapstructure:"level" yaml:"level"`                   // 级别
	Format        string   `mapstructure:"format" yaml:"format"`                 // 输出
	Prefix        string   `mapstructure:"prefix" yaml:"prefix"`                 // 日志前缀
	Director      string   `mapstructure:"director" yaml:"director"`             // 日志文件夹
	ShowLine      bool     `mapstructure:"show-line" yaml:"show-line"`           // 显示行
	EncodeLevel   string   `mapstructure:"encode-level" yaml:"encode-level"`     // 编码级
	StacktraceKey string   `mapstructure:"stacktrace-key" yaml:"stacktrace-key"` // 栈名
	LogInConsole  bool     `mapstructure:"log-in-console" yaml:"log-in-console"` // 输出控制台
	TimeFormat    string   `mapstructure:"time-format" yaml:"time-format"`       // 输出时间格式
	MaxSize       int      `mapstructure:"max-size" yaml:"max-size"`             // 在进行切割之前，日志文件的最大大小（以MB为单位）
	MaxBackups    int      `mapstructure:"max-backups" yaml:"max-backups"`       // 保留旧文件的最大个数
	MaxAge        int      `mapstructure:"max-age" yaml:"max-age"`               // 保留旧文件的最大天数
	Compress      bool     `mapstructure:"compress" yaml:"compress"`             // 是否压缩/归档旧文件
	SkipPaths     []string `mapstructure:"skip-paths" yaml:"skip-paths"`         // 请求时不记录日志的位置
}

type jwt struct {
	JwtTimeout        int    `mapstructure:"jwt-timeout" yaml:"jwt-timeout"` // second
	Secret            string `mapstructure:"secret" yaml:"secret"`           // 加密方式
	DefaultContextKey string `mapstructure:"default-context-key" yaml:"default-context-key"`
}

type decoder struct {
	TemporaryFilePath string `mapstructure:"temporary_file_path" yaml:"temporary_file_path"` // 临时文件存放位置
	ScriptsPath       string `mapstructure:"scripts_path" yaml:"scripts_path"`               // 脚本路径
	DeleteFilePeriod  int    `mapstructure:"delete_file_period" yaml:"delete_file_period"`   // 删除临时文件间隔（小时）
	FileTimeOut       int    `mapstructure:"file_timeout" yaml:"file_timeout"`               // 文件超时时间 （分钟）
}
