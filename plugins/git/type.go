package git

import "time"

type (
	Repo struct {
		Clone string
	}

	Build struct {
		Path   string
		Event  string
		Number int
		Commit string
		Ref    string
	}

	Netrc struct {
		Machine  string
		Login    string
		Password string
	}

	Config struct {
		Depth           int
		Recursive       bool
		SkipVerify      bool
		Tags            bool
		Submodules      map[string]string
		SubmoduleRemote bool
	}

	Backoff struct {
		Attempts int
		Duration time.Duration
	}
)
