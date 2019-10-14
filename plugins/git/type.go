package git

import "time"

type (
	Repo struct {
		Clone string
	}

	Build struct {
		//  仓库地址
		Path   string
		Event  string
		Number int
		Commit string
		Ref    string
	}

	Config struct {
		Depth           int
		Recursive       bool
		SkipVerify      bool
		Tags            bool
		//Submodules      map[string]string
		//SubmoduleRemote bool
	}

	Backoff struct {
		Attempts int
		Duration time.Duration
	}
)
