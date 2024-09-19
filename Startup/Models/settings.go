package Models

type Settings struct {
	Port                  int    `json:"port"`
	ConnectionString      string `json:"connectionString"`
	MaxAllowedConnections int    `json:"maxAllowedConnections"`
	MaxIdleConnection     int    `json:"maxIdleConnection"`
	MaxConnectionTime     int    `json:"maxConnectionTime"`
}
