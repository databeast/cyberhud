package clone

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// RepoConfig holds the configuration for cloning/updating an upstream repository.
type RepoConfig struct {
	URL      string // e.g. "https://github.com/fcambus/spleen.git"
	CacheDir string // e.g. "tools/fontgen/.cache/spleen"
}

// EnsureRepo clones the repository if no cache exists, or pulls updates if it does.
// If the network is unavailable but a cache exists, it prints a warning and continues.
// If the network is unavailable and no cache exists, it returns an error.
func EnsureRepo(cfg RepoConfig) error {
	if cacheExists(cfg.CacheDir) {
		return pullRepo(cfg)
	}
	return cloneRepo(cfg)
}

// cacheExists reports whether the cache directory already exists.
func cacheExists(dir string) bool {
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

// cloneRepo performs a shallow clone of the repository into CacheDir.
// Returns an error if the clone fails (no cache to fall back on).
func cloneRepo(cfg RepoConfig) error {
	cmd := exec.Command("git", "clone", "--depth=1", cfg.URL, cfg.CacheDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("fontgen: cannot clone %s: %s (no cache available)", cfg.URL, firstLine(out, err))
	}
	return nil
}

// pullRepo performs a git pull in the existing cache directory.
// If the pull fails (e.g. network unavailable), it prints a warning to stderr
// and returns nil so the pipeline continues with stale cache.
func pullRepo(cfg RepoConfig) error {
	cmd := exec.Command("git", "-C", cfg.CacheDir, "pull")
	_, err := cmd.CombinedOutput()
	if err != nil {
		repoName := filepath.Base(cfg.CacheDir)
		fmt.Fprintf(os.Stderr, "fontgen: warning: network unavailable, using cached %s\n", repoName)
		return nil
	}
	return nil
}

// firstLine extracts a short description from git command output or the error itself.
func firstLine(output []byte, err error) string {
	s := string(output)
	// Try to use the output as the message if non-empty.
	if len(s) > 0 {
		// Take only the first line for brevity.
		for i, ch := range s {
			if ch == '\n' || ch == '\r' {
				s = s[:i]
				break
			}
		}
		if len(s) > 0 {
			return s
		}
	}
	return err.Error()
}

// DownloadFile downloads a file from url into destPath within cacheDir.
// If the file already exists in cache, it re-downloads (to stay current).
// If the network fails but cache exists, prints a warning and continues.
// If the network fails and no cache exists, returns an error.
func DownloadFile(url, cacheDir, destPath string) error {
	fullPath := filepath.Join(cacheDir, destPath)

	// Ensure cache directory exists.
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("cannot create cache dir: %w", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		// Network failure.
		if fileExists(fullPath) {
			fmt.Fprintf(os.Stderr, "fontgen: warning: download failed, using cached %s\n", destPath)
			return nil
		}
		return fmt.Errorf("fontgen: cannot download %s: %v (no cache available)", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if fileExists(fullPath) {
			fmt.Fprintf(os.Stderr, "fontgen: warning: download returned %d, using cached %s\n", resp.StatusCode, destPath)
			return nil
		}
		return fmt.Errorf("fontgen: cannot download %s: HTTP %d (no cache available)", url, resp.StatusCode)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("cannot write %s: %w", fullPath, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("download %s: %w", url, err)
	}

	return nil
}

// fileExists reports whether the given path exists as a regular file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
