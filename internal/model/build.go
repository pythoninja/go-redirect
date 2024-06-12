package model

type BuildInfo struct {
	AppName   string `json:"-"`
	Version   string `json:"version"`
	VcsTime   string `json:"time"`
	VcsCommit string `json:"commit"`
}
