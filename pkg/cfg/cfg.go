package cfg

type Config struct {
	Port        string `envconfig:"PORT" required:"true"`
	Suffix      string `envconfig:"SUFFIX" required:"true"`
	Bucket      string `envconfig:"BUCKET" required:"true"`
	StorageType string `envconfig:"STORAGE_TYPE" required:"true"`
	ChannelSize int    `envconfig:"CHANNELS_SIZE" required:"true"`
}
