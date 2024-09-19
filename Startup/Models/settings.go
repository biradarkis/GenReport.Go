package Models

type Settings struct {
	Port              int `json:"port"`
	ConnectionStrings struct {
		GenReport string `json:"GenReport"`
	} `json:"connectionStrings"`
}
