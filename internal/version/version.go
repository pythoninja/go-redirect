package version

import (
	"github.com/pythoninja/go-redirect/internal/model"
	"runtime/debug"
)

const appName = "Go-Redirect"

// version will be rewritten on build with actual application version.
var version = "dev"

func GetBuildInfo() *model.BuildInfo {
	var vcsTime, vcsCommit string //nolint:wsl
	var vcsModified bool          //nolint:wsl

	v := version

	if info, ok := debug.ReadBuildInfo(); ok {
		vcsTime, vcsCommit, vcsModified = processBuildSettings(info.Settings)
	}

	if vcsModified {
		v += "-dirty"
	}

	return &model.BuildInfo{
		AppName:   appName,
		Version:   v,
		VcsTime:   vcsTime,
		VcsCommit: vcsCommit,
	}
}

func processBuildSettings(settings []debug.BuildSetting) (vcsTime, vcsCommit string, vcsModified bool) {
	for _, setting := range settings {
		switch setting.Key {
		case "vcs.time":
			vcsTime = setting.Value
		case "vcs.revision":
			vcsCommit = setting.Value
		case "vcs.modified":
			if setting.Value == "true" {
				vcsModified = setting.Value == "true"
			}
		}
	}

	return vcsTime, vcsCommit, vcsModified
}
