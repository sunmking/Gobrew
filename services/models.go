package services

type FormulaInfo struct {
	Name         string    `json:"name"`
	FullName     string    `json:"full_name"`
	Tap          string    `json:"tap"`
	Desc         string    `json:"desc"`
	Homepage     string    `json:"homepage"`
	License      string    `json:"license"`
	StableVer    string    `json:"stable_version"`
	HeadVer      string    `json:"head_version"`
	LinkedKeg    string    `json:"linked_keg"`
	Pinned       bool      `json:"pinned"`
	Installed    []KegInfo `json:"installed"`
	Dependencies []string  `json:"dependencies"`
}

type KegInfo struct {
	Version               string `json:"version"`
	InstalledAsDependency bool   `json:"installed_as_dependency"`
	InstalledOnRequest    bool   `json:"installed_on_request"`
}

type CaskInfo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Tap         string `json:"tap"`
	Desc        string `json:"desc"`
	Homepage    string `json:"homepage"`
	Version     string `json:"version"`
	Installed   string `json:"installed"`
	AutoUpdates bool   `json:"auto_updates"`
	Token       string `json:"token"`
}

type PackageInfoResult struct {
	Type             string    `json:"type"`
	Name             string    `json:"name"`
	FullName         string    `json:"full_name"`
	Tap              string    `json:"tap"`
	Desc             string    `json:"desc"`
	Homepage         string    `json:"homepage"`
	License          string    `json:"license"`
	CurrentVersion   string    `json:"current_version"`
	InstalledVersion string    `json:"installed_version"`
	LinkedKeg        string    `json:"linked_keg"`
	Pinned           bool      `json:"pinned"`
	AutoUpdates      bool      `json:"auto_updates"`
	Token            string    `json:"token"`
	Installed        []KegInfo `json:"installed"`
	Dependencies     []string  `json:"dependencies"`
}

type OutdatedFormula struct {
	Name              string   `json:"name"`
	InstalledVersions []string `json:"installed_versions"`
	CurrentVersion    string   `json:"current_version"`
	Pinned            bool     `json:"pinned"`
	PinnedVersion     string   `json:"pinned_version"`
}

type OutdatedCask struct {
	Name             string `json:"name"`
	InstalledVersion string `json:"installed_version"`
	CurrentVersion   string `json:"current_version"`
}

type BrewServicesResult struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	User     string `json:"user"`
	File     string `json:"file"`
	ExitCode int    `json:"exit_code"`
}

type TapResult struct {
	Name         string `json:"name"`
	Remote       string `json:"remote"`
	CustomRemote bool   `json:"custom_remote"`
}

type TapDetailResult struct {
	Name         string   `json:"name"`
	Remote       string   `json:"remote"`
	CustomRemote bool     `json:"custom_remote"`
	FormulaNames []string `json:"formula_names"`
	CaskTokens   []string `json:"cask_tokens"`
	LastCommit   string   `json:"last_commit"`
	Branch       string   `json:"branch"`
}

type SearchResult struct {
	Formulae []SearchItem `json:"formulae"`
	Casks    []SearchItem `json:"casks"`
}

type SearchItem struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Desc     string `json:"desc"`
	Tap      string `json:"tap"`
}

type BrewError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *BrewError) Error() string {
	if e == nil {
		return ""
	}
	if e.Details == "" {
		return e.Message
	}
	return e.Message + ": " + e.Details
}

type CommandResult struct {
	Success  bool   `json:"success"`
	Output   string `json:"output"`
	Duration string `json:"duration"`
	Error    string `json:"error,omitempty"`
}

type CleanupResult struct {
	CleanedCount  int    `json:"cleaned_count"`
	SizeReclaimed string `json:"size_reclaimed"`
	Output        string `json:"output"`
}

type AppConfig struct {
	Language             string `json:"language"`
	Theme                string `json:"theme"`
	BrewPath             string `json:"brew_path"`
	BrewfilePath         string `json:"brewfile_path"`
	AutoUpdateInterval   string `json:"auto_update_interval"`
	LogMaxLines          int    `json:"log_max_lines"`
	CheckUpdatesOnLaunch bool   `json:"check_updates_on_launch"`
}

type BrewPathValidation struct {
	Path    string `json:"path"`
	Valid   bool   `json:"valid"`
	Version string `json:"version"`
	Error   string `json:"error"`
}

type BundleCheckResult struct {
	Satisfied bool     `json:"satisfied"`
	Missing   []string `json:"missing"`
	Output    string   `json:"output"`
}

type InstalledListResult struct {
	Formulae []FormulaInfo `json:"formulae"`
	Casks    []CaskInfo    `json:"casks"`
}

type OutdatedListResult struct {
	Formulae []OutdatedFormula `json:"formulae"`
	Casks    []OutdatedCask    `json:"casks"`
}
