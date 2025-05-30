package module

import (
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type SymlinkedDirectory struct {
	/**
	 * Matches the casing returned by `realpath`.  Used to compute the `realpath` of children.
	 * Always has trailing directory separator
	 */
	Real string
	/**
	 * toPath(real).  Stored to avoid repeated recomputation.
	 * Always has trailing directory separator
	 */
	RealPath tspath.Path
}

type SymlinkCache struct {
	symlinkedDirectories           map[tspath.Path]*SymlinkedDirectory
	symlinkedDirectoriesByRealpath core.MultiMap[tspath.Path, string]
	symlinkedFiles                 map[tspath.Path]string
}

/** Gets a map from symlink to realpath. Keys have trailing directory separators. */
func (cache *SymlinkCache) SymlinkedDirectories() map[tspath.Path]*SymlinkedDirectory {
	return cache.symlinkedDirectories
}

/** Gets a map from symlink to realpath */
func (cache *SymlinkCache) SymlinkedFiles() map[tspath.Path]string {
	return cache.symlinkedFiles
}

func (cache *SymlinkCache) SetSymlinkedDirectory(symlink string, symlinkPath tspath.Path, real *SymlinkedDirectory) {
	// Large, interconnected dependency graphs in pnpm will have a huge number of symlinks
	// where both the realpath and the symlink path are inside node_modules/.pnpm. Since
	// this path is never a candidate for a module specifier, we can ignore it entirely.

	// !!! sheetal - all callers check this !containsIgnoredPath(symlinkPath) ?
	if real != nil {
		if _, ok := cache.symlinkedDirectories[symlinkPath]; !ok {
			cache.symlinkedDirectoriesByRealpath.Add(real.RealPath, symlink)
		}
	}
	if cache.symlinkedDirectories == nil {
		cache.symlinkedDirectories = make(map[tspath.Path]*SymlinkedDirectory)
	}
	cache.symlinkedDirectories[symlinkPath] = real
}

func (cache *SymlinkCache) SetSymlinkedFile(symlinkPath tspath.Path, real string) {
	if cache.symlinkedFiles == nil {
		cache.symlinkedFiles = make(map[tspath.Path]string)
	}
	cache.symlinkedFiles[symlinkPath] = real
}
