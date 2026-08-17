package source

import (
	"os/exec"
	"strings"
	"sync"
)

const defaultDesiredTarget = "multi-user.target"

type probeFuncs struct {
	desiredTarget    func() string
	loadingTarget    func() string
	reachedMultiUser func() bool
	listTargets      func() (active int, total int)
}

var probes = probeFuncs{
	desiredTarget:    queryDesiredTarget,
	loadingTarget:    queryLoadingTarget,
	reachedMultiUser: queryReachedMultiUser,
	listTargets:      queryTargetCounts,
}

var state struct {
	sync.Mutex
	desired string
	loading string
	loaded  string
}

var systemctlExec = func(args ...string) (string, error) {
	cmd := exec.Command("systemctl", args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// ReachedMultiUser reports whether multi-user.target is currently active.
func ReachedMultiUser() bool {
	return PollSnapshot().BootComplete
}

func PollSnapshot() Snapshot {
	desired := sanitizeTarget(probes.desiredTarget())
	if desired == "" {
		desired = defaultDesiredTarget
	}
	loading := sanitizeTarget(probes.loadingTarget())
	bootComplete := probes.reachedMultiUser()
	active, total := probes.listTargets()

	state.Lock()
	defer state.Unlock()

	state.desired = desired
	if bootComplete {
		return Snapshot{
			Desired:       state.desired,
			Loading:       state.loading,
			Loaded:        state.loaded,
			BootComplete:  true,
			ActiveTargets: active,
			TotalTargets:  total,
		}
	}
	if loading == "" {
		loading = state.desired
	}
	if state.loading != "" && state.loading != loading {
		state.loaded = state.loading
	}
	state.loading = loading
	return Snapshot{
		Desired:       state.desired,
		Loading:       state.loading,
		Loaded:        state.loaded,
		ActiveTargets: active,
		TotalTargets:  total,
	}
}

func queryDesiredTarget() string {
	out, err := systemctlExec("get-default", "--no-pager")
	if err != nil {
		return ""
	}
	return out
}

func queryLoadingTarget() string {
	out, err := systemctlExec("list-jobs", "--type=target", "--no-legend", "--no-pager")
	if err != nil || strings.TrimSpace(out) == "" {
		return ""
	}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		target := sanitizeTarget(fields[0])
		if target != "" {
			return target
		}
	}
	return ""
}

func queryReachedMultiUser() bool {
	_, err := systemctlExec("is-active", "--quiet", "multi-user.target", "--no-pager")
	return err == nil
}

// queryTargetCounts returns the number of active targets and total known targets.
// It uses `systemctl list-units --type=target --no-legend --no-pager --all` to
// enumerate all targets and counts how many have an "active" state.
func queryTargetCounts() (active int, total int) {
	out, err := systemctlExec("list-units", "--type=target", "--no-legend", "--no-pager", "--all")
	if err != nil || strings.TrimSpace(out) == "" {
		return 0, 0
	}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		// Fields: UNIT LOAD ACTIVE SUB DESCRIPTION...
		// We care about the ACTIVE column (index 2).
		total++
		if fields[2] == "active" {
			active++
		}
	}
	return active, total
}

func sanitizeTarget(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return ""
	}
	target := strings.TrimSpace(fields[0])
	if strings.HasSuffix(target, ".target") {
		return target
	}
	return ""
}

// SanitizeTarget exposes target normalization for package-level tests.
func SanitizeTarget(value string) string { return sanitizeTarget(value) }

// ResetStateForTest clears the cached transition state used while polling.
func ResetStateForTest() {
	state.Lock()
	state.desired = ""
	state.loading = ""
	state.loaded = ""
	state.Unlock()
}

// SetProbesForTest replaces systemctl probe functions for deterministic tests.
func SetProbesForTest(desired func() string, loading func() string, reached func() bool, targets func() (int, int)) {
	probes = probeFuncs{
		desiredTarget:    desired,
		loadingTarget:    loading,
		reachedMultiUser: reached,
		listTargets:      targets,
	}
}
