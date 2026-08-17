package clone

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// setupBareRepo creates a temporary bare git repository with one commit,
// suitable for use as a remote URL in tests.
func setupBareRepo(t *testing.T) string {
	t.Helper()

	// Create a temporary directory for the "remote" bare repo.
	bareDir := t.TempDir()

	// Initialize a bare repository.
	run(t, bareDir, "git", "init", "--bare")

	// We need a working clone to push an initial commit into the bare repo.
	workDir := t.TempDir()
	run(t, workDir, "git", "clone", bareDir, ".")
	run(t, workDir, "git", "config", "user.email", "test@test.com")
	run(t, workDir, "git", "config", "user.name", "Test")

	// Create a file and commit it.
	if err := os.WriteFile(filepath.Join(workDir, "README.md"), []byte("# test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run(t, workDir, "git", "add", ".")
	run(t, workDir, "git", "commit", "-m", "initial commit")
	run(t, workDir, "git", "push", "origin", "HEAD")

	return bareDir
}

// run executes a git command in the given directory and fails the test on error.
func run(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed in %s: %v\n%s", name, args, dir, err, out)
	}
}

// TestEnsureRepoClonesFresh verifies that EnsureRepo clones a repository
// into an empty cache directory.
func TestEnsureRepoClonesFresh(t *testing.T) {
	bareRepo := setupBareRepo(t)
	cacheDir := filepath.Join(t.TempDir(), "cache", "repo")

	cfg := RepoConfig{
		URL:      bareRepo,
		CacheDir: cacheDir,
	}

	err := EnsureRepo(cfg)
	if err != nil {
		t.Fatalf("EnsureRepo failed on fresh clone: %v", err)
	}

	// Verify the cache directory was created and contains a git repo.
	if _, err := os.Stat(filepath.Join(cacheDir, ".git")); os.IsNotExist(err) {
		t.Fatal("cache directory does not contain a .git directory after clone")
	}

	// Verify the cloned content is present.
	readme := filepath.Join(cacheDir, "README.md")
	if _, err := os.Stat(readme); os.IsNotExist(err) {
		t.Fatal("cloned repo does not contain README.md")
	}
}

// TestEnsureRepoPullsExisting verifies that when a cache already exists,
// EnsureRepo performs a pull (not a re-clone). We verify this by adding
// a new commit to the bare repo and checking that EnsureRepo picks it up.
func TestEnsureRepoPullsExisting(t *testing.T) {
	bareRepo := setupBareRepo(t)
	cacheDir := filepath.Join(t.TempDir(), "cache", "repo")

	cfg := RepoConfig{
		URL:      bareRepo,
		CacheDir: cacheDir,
	}

	// First call: clone.
	if err := EnsureRepo(cfg); err != nil {
		t.Fatalf("initial EnsureRepo failed: %v", err)
	}

	// Push a new commit to the bare repo.
	workDir := t.TempDir()
	run(t, workDir, "git", "clone", bareRepo, ".")
	run(t, workDir, "git", "config", "user.email", "test@test.com")
	run(t, workDir, "git", "config", "user.name", "Test")
	if err := os.WriteFile(filepath.Join(workDir, "CHANGELOG.md"), []byte("# changes\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run(t, workDir, "git", "add", ".")
	run(t, workDir, "git", "commit", "-m", "add changelog")
	run(t, workDir, "git", "push", "origin", "HEAD")

	// Second call: should pull, not re-clone.
	if err := EnsureRepo(cfg); err != nil {
		t.Fatalf("second EnsureRepo failed: %v", err)
	}

	// Verify the new file is now in the cache (pulled from remote).
	changelog := filepath.Join(cacheDir, "CHANGELOG.md")
	if _, err := os.Stat(changelog); os.IsNotExist(err) {
		t.Fatal("pull did not fetch new CHANGELOG.md; EnsureRepo may have re-cloned or skipped pull")
	}
}

// TestEnsureRepoFailsNoCache verifies that when the URL is invalid and
// no cache exists, EnsureRepo returns an error with a descriptive message.
func TestEnsureRepoFailsNoCache(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "nonexistent", "cache")

	cfg := RepoConfig{
		URL:      "https://invalid.example.com/nonexistent/repo.git",
		CacheDir: cacheDir,
	}

	err := EnsureRepo(cfg)
	if err == nil {
		t.Fatal("expected error when network unavailable and no cache, got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "no cache available") {
		t.Errorf("error message should mention 'no cache available', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, cfg.URL) {
		t.Errorf("error message should contain the URL, got: %s", errMsg)
	}
}

// TestEnsureRepoWarnsOnNetworkFailure verifies that when a cache exists
// but the network is unavailable (invalid URL), EnsureRepo succeeds and
// prints a warning to stderr about using the cached clone.
func TestEnsureRepoWarnsOnNetworkFailure(t *testing.T) {
	// First, create a valid clone to serve as existing cache.
	bareRepo := setupBareRepo(t)
	cacheDir := filepath.Join(t.TempDir(), "cache", "repo")

	cfg := RepoConfig{
		URL:      bareRepo,
		CacheDir: cacheDir,
	}

	// Clone successfully first.
	if err := EnsureRepo(cfg); err != nil {
		t.Fatalf("initial clone failed: %v", err)
	}

	// Now change the git remote in the cached clone to an invalid URL
	// to simulate network failure on the next pull.
	cfg.URL = "https://invalid.example.com/unreachable/repo.git"
	run(t, cacheDir, "git", "remote", "set-url", "origin", cfg.URL)

	// Capture stderr to verify the warning is printed.
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w

	// Call EnsureRepo with invalid URL but existing cache.
	ensureErr := EnsureRepo(cfg)

	// Restore stderr and read captured output.
	w.Close()
	os.Stderr = oldStderr
	var buf bytes.Buffer
	buf.ReadFrom(r)
	stderrOutput := buf.String()

	// EnsureRepo should succeed (return nil) since cache exists.
	if ensureErr != nil {
		t.Fatalf("expected nil error when cache exists but network fails, got: %v", ensureErr)
	}

	// Verify warning was printed to stderr.
	if !strings.Contains(stderrOutput, "warning") {
		t.Errorf("expected warning on stderr, got: %q", stderrOutput)
	}
	if !strings.Contains(stderrOutput, "cached") {
		t.Errorf("expected stderr warning to mention 'cached', got: %q", stderrOutput)
	}
}
