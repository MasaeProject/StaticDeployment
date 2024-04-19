package main

type Project struct {
	Name       string        `json:"name"`
	Source     string        `json:"source"`
	Additional *[]string     `json:"additional"`
	Run        *Run          `json:"run"`
	PreRun     *Run          `json:"prerun"`
	Replace    []ReplaceItem `json:"replace"`
}

type ReplaceItem struct {
	Name       string          `json:"to"`
	Replace    []ReplaceDetail `json:"replace"`
	Additional *[]string       `json:"additional"`
	Run        *Run            `json:"run"`
	PreRun     *Run            `json:"prerun"`
}

type Run struct {
	Default *[][]string `json:"default"`
	Windows *[][]string `json:"windows"`
	Linux   *[][]string `json:"linux"`
	Darwin  *[][]string `json:"darwin"`
}

type ReplaceDetail struct {
	Old string `json:"old"`
	New string `json:"new"`
	Num int    `json:"num"`
}

type FileData struct {
	Path       string
	LoadString string
	LoadSize   int
	NewString  string
	NewSize    int
}

type Names struct {
	Project string
	Replace string
}
