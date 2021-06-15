package webhook

import (
	"time"
)

// EventType describes the type of the incoming event
type EventType int

// All supported event types are defined below
const (
	EventTypePush = iota
	EventTypePullRequest
)

var eventTypeTrans = map[string]EventType{
	"push":         EventTypePush,
	"pull_request": EventTypePullRequest,
}

// String returns the event type as a string
func (typ EventType) String() string {
	switch typ {
	case EventTypePush:
		return "push"
	case EventTypePullRequest:
		return "pull request"
	default:
		return "unknown"
	}
}

// Event describes a Gitea webhook event
type Event struct {
	Secret      string      `json:"secret"`
	Action      string      `json:"action"`
	Number      int         `json:"number"`
	Ref         string      `json:"ref"`
	Before      string      `json:"before"`
	After       string      `json:"after"`
	CompareURL  string      `json:"compare_url"`
	Commits     []Commit    `json:"commits"`
	PullRequest PullRequest `json:"pull_reuqest"`
	Repository  Repository  `json:"repository"`
	Pusher      User        `json:"pusher"`
	Sender      User        `json:"sender"`
}

// Commit describes a single commit
type Commit struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	URL       string    `json:"url"`
	Author    GitUser   `json:"author"`
	Committer GitUser   `json:"committer"`
	Timestamp time.Time `json:"timestamp"`
}

// PullRequest describes a open pull request
type PullRequest struct {
	ID             int        `json:"id"`
	URL            string     `json:"url"`
	Number         int        `json:"number"`
	User           User       `json:"user"`
	Title          string     `json:"title"`
	Body           string     `json:"body"`
	Labels         []string   `json:"labels"`
	Assignee       User       `json:"assignee"`
	Assignees      []User     `json:"assignees"`
	State          string     `json:"state"`
	Comments       int        `json:"comments"`
	HtmlURL        string     `json:"html_url"`
	DiffURL        string     `json:"diff_url"`
	PatchURL       string     `json:"patch_url"`
	Mergeable      bool       `json:"mergeable"`
	Merged         bool       `json:"mergeable"`
	MergedAt       *time.Time `json:"merged_at"`
	MergeCommitSHA *string    `json:"merge_commit_sha"`
	MergedBy       *User      `json:"merged_by"`
	Base           Ref        `json:"base"`
	Head           Ref        `json:"head"`
	MergeBase      string     `json:"merge_base"`
	DueDate        *time.Time `json:"due_date"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	ClosedAt       *time.Time `json:"closed_at"`
}

// Ref describes a reference to a specific commit in a repository
type Ref struct {
	Label  string     `json:"label"`
	Ref    string     `json:"ref"`
	SHA    string     `json:"sha"`
	RepoID string     `json:"repo_id"`
	Repo   Repository `json:"repo"`
}

// Repository describes a Gitea git repository
type Repository struct {
	ID                        int                       `json:"id"`
	Owner                     User                      `json:"owner"`
	Name                      string                    `json:"name"`
	FullName                  string                    `json:"full_name"`
	Description               string                    `json:"description"`
	Empty                     bool                      `json:"empty"`
	Private                   bool                      `json:"private"`
	Fork                      bool                      `json:"fork"`
	Parent                    *Repository               `json:"parent"`
	Mirror                    bool                      `json:"mirror"`
	Size                      int                       `json:"size"`
	HtmlURL                   string                    `json:"html_url"`
	SshURL                    string                    `json:"ssh_url"`
	CloneURL                  string                    `json:"clone_url"`
	OriginalURL               string                    `json:"original_url"`
	Website                   string                    `json:"website"`
	StarsCount                int                       `json:"stars_count"`
	ForksCount                int                       `json:"forks_count"`
	WatchersCount             int                       `json:"watchers_count"`
	OpenIssuesCount           int                       `json:"open_issues_count"`
	DefaultBranch             string                    `json:"default_branch"`
	Archived                  bool                      `json:"archived"`
	CreatedAt                 time.Time                 `json:"created_at"`
	UpdatedAt                 time.Time                 `json:"updated_at"`
	Permissions               RepositoryPermissions     `json:"permissions"`
	HasIssues                 bool                      `json:"has_issues"`
	InternalTracker           RepositoryInternalTracker `json:"internal_tracker"`
	HasWiki                   bool                      `json:"has_wiki"`
	HasPullRequests           bool                      `json:"has_pull_requests"`
	IgnoreWhitespaceConflicts bool                      `json:"ignore_whitespace_conflicts"`
	AllowMergeCommits         bool                      `json:"allow_merge_commits"`
	AllowRebase               bool                      `json:"allow_rebase"`
	AllowRebaseExplicit       bool                      `json:"allow_rebase_explicit"`
	AllowSquashMerge          bool                      `json:"allow_squash_merge"`
	AvatarURL                 string                    `json:"avatar_url"`
}

type RepositoryPermissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

type RepositoryInternalTracker struct {
	EnableTimeTracker                bool `json:"enable_time_tracker"`
	AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
	EnableIssueDependencies          bool `json:"enable_issue_dependencies"`
}

// User describes a Gitea user
type User struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	Language  string    `json:"language"`
	IsAdmin   bool      `json:"is_admin"`
	LastLogin time.Time `json:"last_login"`
	Created   time.Time `json:"created"`
	Username  string    `json:"username"`
}

// GitUser describes a git user (in commits)
type GitUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// CommitStatusState is used to set/get the current state of a single commit
type CommitStatusState string

const (
	// CommitStatusPending is for when the CommitStatus is Pending
	CommitStatusPending CommitStatusState = "pending"
	// CommitStatusSuccess is for when the CommitStatus is Success
	CommitStatusSuccess CommitStatusState = "success"
	// CommitStatusError is for when the CommitStatus is Error
	CommitStatusError CommitStatusState = "error"
	// CommitStatusFailure is for when the CommitStatus is Failure
	CommitStatusFailure CommitStatusState = "failure"
	// CommitStatusWarning is for when the CommitStatus is Warning
	CommitStatusWarning CommitStatusState = "warning"
)

// CommitStatus describes the state of a single commit
type CommitStatus struct {
	ID          int64
	Index       int64
	RepoID      int64
	Repo        *Repository
	State       CommitStatusState
	SHA         string
	TargetURL   string
	Description string
	ContextHash string
	Context     string
	Creator     User
	CreatorID   int64
}

// CreateStatusOption is used to update the status flag for a single commit
type CreateStatusOption struct {
	Context     string            `json:"context"`
	Description string            `json:"description"`
	State       CommitStatusState `json:"state"`
	TargetURL   string            `json:"target_url"` // Used to create a link from gitea to another system
}
