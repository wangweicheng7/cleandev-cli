package types

type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

type Profile string

const (
	ProfileSafe       Profile = "safe"
	ProfileDev        Profile = "dev"
	ProfileAggressive Profile = "aggressive"
)

type Rule struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	Paths          []string  `json:"paths"`
	RiskLevel      RiskLevel `json:"risk_level"`
	Profiles       []Profile `json:"profiles"`
	DeleteStrategy string    `json:"delete_strategy"`
}

type Candidate struct {
	RuleID          string    `json:"rule_id"`
	Name            string    `json:"name"`
	Category        string    `json:"category"`
	Path            string    `json:"path"`
	Size            int64     `json:"size"`
	Reclaimable     int64     `json:"reclaimable_bytes"`
	RiskLevel       RiskLevel `json:"risk_level"`
	Status          string    `json:"status"`
	ReportOnly      bool      `json:"report_only"`
	Reason          string    `json:"reason"`
	DeleteStrategy  string    `json:"delete_strategy"`
	LastAccessOrMod string    `json:"last_access"`
}

type Action struct {
	Path        string    `json:"path"`
	Category    string    `json:"category"`
	RiskLevel   RiskLevel `json:"risk_level"`
	Reclaimable int64     `json:"reclaimable_bytes"`
	Action      string    `json:"action"`
	Reason      string    `json:"reason"`
}

type PlanSummary struct {
	TotalItems       int                 `json:"total_items"`
	ReclaimableBytes int64               `json:"reclaimable_bytes"`
	ByCategory       map[string]int64    `json:"by_category"`
	ByRisk           map[RiskLevel]int64 `json:"by_risk"`
}

type Plan struct {
	Summary PlanSummary `json:"summary"`
	Actions []Action    `json:"actions"`
}

type Config struct {
	Command        string
	Profile        Profile
	JSONOutput     bool
	Confirm        bool
	Categories     map[string]struct{}
	ConfigPath     string
	IncludePaths   []string
	ExcludePaths   []string
	ProtectedPaths []string
	Rules          []Rule
	AuditLogDir    string
}
