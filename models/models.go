package models

type Token struct {
    Addr    string  `json:"addr"`
    Port    int     `json:"port"`
    Host    string  `json:"host"`
    Key     string  `json:"key"`
}
