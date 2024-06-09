package model

type Config struct {
	HTTP     Http     `json:"http"`
	Postgres Postgres `json:"postgres"`
	JWT      JWT      `json:"jwt"`
	Email    Email    `json:"email"`
}

type Http struct {
	Port string `json:"port"`
}

type Postgres struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	SSL       string `json:"ssl"`
	Migration string `json:"migration"`
}

type Email struct {
	From string `json:"from"`
	To   string `json:"to"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Port string `json:"port"`
}

type JWT struct {
	Secret string `json:"secret"`
}
