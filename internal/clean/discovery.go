package clean

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DiscoverOptions struct {
	Enabled  bool
	Roots    []string
	MaxDepth int
	Refresh  bool
	CacheTTL time.Duration
}

type discoveredProject struct {
	Root string
	Name string
}

type discoveryCache struct {
	GeneratedAt time.Time          `json:"generated_at"`
	Roots       []string           `json:"roots"`
	MaxDepth    int                `json:"max_depth"`
	Projects    []discoveredProject `json:"projects"`
}

func discoverProjects(ctx context.Context, home string, opts DiscoverOptions) ([]discoveredProject, error) {
	if !opts.Enabled {
		return nil, nil
	}
	maxDepth := opts.MaxDepth
	if maxDepth <= 0 {
		maxDepth = 4
	}

	roots := normalizeDiscoverRoots(home, opts.Roots)
	ttl := opts.CacheTTL
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	if !opts.Refresh {
		if cached, ok := loadDiscoveryCache(home, roots, maxDepth, ttl); ok {
			// prune missing paths so cache heals itself when projects disappear
			var alive []discoveredProject
			for _, p := range cached {
				if st, err := os.Stat(p.Root); err == nil && st.IsDir() {
					alive = append(alive, p)
				}
			}
			_ = saveDiscoveryCache(home, roots, maxDepth, alive)
			return alive, nil
		}
	}

	seen := map[string]bool{}
	var out []discoveredProject

	for _, root := range roots {
		info, err := os.Stat(root)
		if err != nil || !info.IsDir() {
			continue
		}

		err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			if !d.IsDir() {
				return nil
			}
			if strings.HasPrefix(d.Name(), ".") && d.Name() != ".dart_tool" && d.Name() != ".gradle" {
				return filepath.SkipDir
			}

			rel, _ := filepath.Rel(root, path)
			depth := 0
			if rel != "." {
				depth = strings.Count(filepath.ToSlash(rel), "/") + 1
			}
			if depth > maxDepth {
				return filepath.SkipDir
			}

			if isProjectDir(path) {
				cleaned := filepath.Clean(path)
				if !seen[cleaned] {
					seen[cleaned] = true
					out = append(out, discoveredProject{
						Root: cleaned,
						Name: filepath.Base(cleaned),
					})
				}
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	_ = saveDiscoveryCache(home, roots, maxDepth, out)
	return out, nil
}

func normalizeDiscoverRoots(home string, roots []string) []string {
	if len(roots) == 0 {
		roots = []string{"~/Code", "~/Projects", "~/workspace"}
	}
	var out []string
	for _, r := range roots {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		if strings.HasPrefix(r, "~/") {
			r = filepath.Join(home, strings.TrimPrefix(r, "~/"))
		}
		if !filepath.IsAbs(r) {
			r = filepath.Join(home, r)
		}
		out = append(out, filepath.Clean(r))
	}
	return out
}

func isProjectDir(dir string) bool {
	markers := []string{
		"package.json",
		"pubspec.yaml",
		"go.mod",
		"pom.xml",
		"build.gradle",
		"build.gradle.kts",
		"settings.gradle",
		"settings.gradle.kts",
		"ios",
		"android",
	}
	for _, m := range markers {
		if _, err := os.Stat(filepath.Join(dir, m)); err == nil {
			return true
		}
	}
	return false
}

func projectTargets(p discoveredProject) []Item {
	inProject := func(rel string) string { return filepath.Join(p.Root, rel) }

	return []Item{
		{
			ID:         "proj-node_modules:" + p.Root,
			Name:       p.Name + " node_modules",
			Path:       inProject("node_modules"),
			Category:   CategoryBuild,
			ProfileMin: ProfileDev,
			Mode:       ModeDelete,
			Reason:     "dependency install output, can be reinstalled",
		},
		{
			ID:         "proj-dist:" + p.Root,
			Name:       p.Name + " dist",
			Path:       inProject("dist"),
			Category:   CategoryBuild,
			ProfileMin: ProfileDev,
			Mode:       ModeDelete,
			Reason:     "build output, can be regenerated",
		},
		{
			ID:         "proj-build:" + p.Root,
			Name:       p.Name + " build",
			Path:       inProject("build"),
			Category:   CategoryBuild,
			ProfileMin: ProfileDev,
			Mode:       ModeDelete,
			Reason:     "build output, can be regenerated",
		},
		{
			ID:         "proj-dart-tool:" + p.Root,
			Name:       p.Name + " .dart_tool",
			Path:       inProject(".dart_tool"),
			Category:   CategoryBuild,
			ProfileMin: ProfileDev,
			Mode:       ModeDelete,
			Reason:     "Flutter/Dart tool cache, can be regenerated",
		},
		{
			ID:         "proj-target:" + p.Root,
			Name:       p.Name + " target",
			Path:       inProject("target"),
			Category:   CategoryBuild,
			ProfileMin: ProfileDev,
			Mode:       ModeDelete,
			Reason:     "Rust/Java build output, can be regenerated",
		},
		{
			ID:         "proj-out:" + p.Root,
			Name:       p.Name + " out",
			Path:       inProject("out"),
			Category:   CategoryBuild,
			ProfileMin: ProfileDev,
			Mode:       ModeDelete,
			Reason:     "build output, can be regenerated",
		},
		{
			ID:         "proj-android-gradle:" + p.Root,
			Name:       p.Name + " android/.gradle",
			Path:       inProject("android/.gradle"),
			Category:   CategoryBuild,
			ProfileMin: ProfileDev,
			Mode:       ModeDelete,
			Reason:     "Android project-local gradle cache",
		},
		{
			ID:         "proj-ios-pods:" + p.Root,
			Name:       p.Name + " ios/Pods",
			Path:       inProject("ios/Pods"),
			Category:   CategoryBuild,
			ProfileMin: ProfileDev,
			Mode:       ModeDelete,
			Reason:     "CocoaPods install output, can be reinstalled",
		},
	}
}

func loadDiscoveryCache(home string, roots []string, maxDepth int, ttl time.Duration) ([]discoveredProject, bool) {
	b, err := os.ReadFile(discoveryCachePath(home))
	if err != nil {
		return nil, false
	}
	var c discoveryCache
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, false
	}
	if c.MaxDepth != maxDepth {
		return nil, false
	}
	if !sameStringSlice(c.Roots, roots) {
		return nil, false
	}
	if time.Since(c.GeneratedAt) > ttl {
		return nil, false
	}
	return c.Projects, true
}

func saveDiscoveryCache(home string, roots []string, maxDepth int, projects []discoveredProject) error {
	if projects == nil {
		projects = []discoveredProject{}
	}
	c := discoveryCache{
		GeneratedAt: time.Now(),
		Roots:       roots,
		MaxDepth:    maxDepth,
		Projects:    projects,
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	p := discoveryCachePath(home)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, b, 0o644)
}

func discoveryCachePath(home string) string {
	return filepath.Join(home, ".cache", "devclean", "discovery-cache.json")
}

func sameStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

