package main

type Project struct {
	Name       string        `json:"name"`
	Source     string        `json:"source"`
	Additional []string      `json:"additional"`
	Exec       [][]string    `json:"exec"`
	Replace    []ReplaceItem `json:"replace"`
}

type ReplaceItem struct {
	To         string          `json:"to"`
	Replace    []ReplaceDetail `json:"replace"`
	Additional []string        `json:"additional"`
	Exec       [][]string      `json:"exec"`
}

type ReplaceDetail struct {
	Old string `json:"old"`
	New string `json:"new"`
	Num int    `json:"num"`
}

type FileData struct {
	LoadString string
	LoadSize   int
	NewString  string
	NewSize    int
}
