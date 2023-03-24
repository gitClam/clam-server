package config

type config struct {
	System system `mapstructure:"system" yaml:"system"`
	Zap    zap    `mapstructure:"zap" yaml:"zap"`
}

type system struct {
	Host string `yaml:"host"` // 端口号
}

type zap struct {
	Level         string `yaml:"level"`          // 级别
	Format        string `yaml:"format"`         // 输出
	Prefix        string `yaml:"prefix"`         // 日志前缀
	Director      string `yaml:"director"`       // 日志文件夹
	ShowLine      bool   `yaml:"showLine"`       // 显示行
	EncodeLevel   string `yaml:"encode-level"`   // 编码级
	StacktraceKey string `yaml:"stacktrace-key"` // 栈名
	LogInConsole  bool   `yaml:"log-in-console"` // 输出控制台
	TimeFormat    string `yaml:"timeFormat"`     // 输出时间格式
	MaxSize       int    `yaml:"max-size"`       // 在进行切割之前，日志文件的最大大小（以MB为单位）
	MaxBackups    int    `yaml:"max-backups"`    // 保留旧文件的最大个数
	MaxAge        int    `yaml:"max-age"`        // 保留旧文件的最大天数
	Compress      bool   `yaml:"compress"`       // 是否压缩/归档旧文件
}
