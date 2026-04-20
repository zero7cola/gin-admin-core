package setting

type Setting struct {
	App     AppConfig     `mapstructure:"app" yaml:"app"`
	Storage StorageConfig `mapstructure:"storage" yaml:"storage"`
	JWT     JWTConfig     `mapstructure:"jwt" yaml:"jwt"`
	Captcha CaptchaConfig `mapstructure:"captcha" yaml:"captcha"`
	Paging  PagingConfig  `mapstructure:"paging" yaml:"paging"`
}

var GlobalSetting *Setting

type AppConfig struct {
	Name     string `mapstructure:"name" yaml:"name"`
	Key      string `mapstructure:"key" yaml:"key"`
	Url      string `mapstructure:"url" yaml:"url"`
	HttpPort string `mapstructure:"http_port" yaml:"http_port"`
	FileUrl  string `mapstructure:"file_url" yaml:"file_url"`
	Env      string `mapstructure:"env" yaml:"env"`
	Version  string `mapstructure:"version" yaml:"version"`
	Debug    bool   `mapstructure:"debug" yaml:"debug"`
	Timezone string `mapstructure:"timezone" yaml:"timezone"`
}

type StorageConfig struct {
	Driver    string              `mapstructure:"driver" yaml:"driver"`
	SizeLimit int64               `mapstructure:"size_limit" yaml:"size_limit"`
	Ext       []string            `mapstructure:"ext" yaml:"ext"`
	Local     *LocalStorageConfig `mapstructure:"local" yaml:"local"`
	Oss       *OssStorageConfig   `mapstructure:"oss" yaml:"oss"`
}

type LocalStorageConfig struct {
	Path         string `mapstructure:"path" yaml:"path"`
	Domain       string `mapstructure:"domain" yaml:"domain"`
	StaticPrefix string `mapstructure:"static" yaml:"static"`
}

type OssStorageConfig struct {
	KeyId     string `mapstructure:"key_id" yaml:"key_id"`
	KeySecret string `mapstructure:"key_secret" yaml:"key_secret"`
	Region    string `mapstructure:"region" yaml:"region"`
	Bucket    string `mapstructure:"bucket" yaml:"bucket"`
	Domain    string `mapstructure:"domain" yaml:"domain"`
}

type PagingConfig struct {
	PerPage         int    `mapstructure:"perpage" yaml:"perpage"`
	UrlQueryOrder   string `mapstructure:"url_query_order" yaml:"url_query_order"`
	UrlQuerySort    string `mapstructure:"url_query_sort" yaml:"url_query_sort"`
	UrlQueryPage    string `mapstructure:"url_query_page" yaml:"url_query_page"`
	UrlQueryPerPage string `mapstructure:"url_query_per_page" yaml:"url_query_per_page"`
}

type CaptchaConfig struct {
	Height     int     `mapstructure:"height" yaml:"height"`
	Width      int     `mapstructure:"width" yaml:"width"`
	Length     int     `mapstructure:"length" yaml:"length"`
	Maxskew    float64 `mapstructure:"maxskew" yaml:"maxskew"`
	Dotcount   int     `mapstructure:"dotcount" yaml:"dotcount"`
	ExpireTime int     `mapstructure:"expire_time" yaml:"expire_time"`
}

type JWTConfig struct {
	ExpireTime     int `mapstructure:"expire_time" yaml:"expire_time"`           // 过期时间，单位是分钟，一般不超过两个小时
	MaxReFreshTime int `mapstructure:"max_refresh_time" yaml:"max_refresh_time"` // 允许刷新时间，单位分钟，从 Token 的签名时间算起
}
