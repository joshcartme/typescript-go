package compiler

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type ReferenceResolutionHost struct {
	host              CompilerHost
	outputDtsToSource map[tspath.Path]*tsoptions.SourceAndProjectReference
	dtsDirectories    core.Set[tspath.Path]
	symlinkCache      module.SymlinkCache
}

var _ module.ResolutionHost = (*ReferenceResolutionHost)(nil)

func (h *ReferenceResolutionHost) FS() vfs.FS {
	return h
}

func (h *ReferenceResolutionHost) GetCurrentDirectory() string {
	return h.host.GetCurrentDirectory()
}

func (h *ReferenceResolutionHost) Trace(msg string) {
	h.host.Trace(msg)
}

// UseCaseSensitiveFileNames returns true if the file system is case-sensitive.
func (h *ReferenceResolutionHost) UseCaseSensitiveFileNames() bool {
	return h.host.FS().UseCaseSensitiveFileNames()
}

// FileExists returns true if the file exists.
func (h *ReferenceResolutionHost) FileExists(path string) bool {
	if h.host.FS().FileExists(path) {
		return true
	}
	if !tspath.IsDeclarationFileName(path) {
		return false
	}
	// Project references go to source file instead of .d.ts file
	return h.fileOrDirectoryExistsUsingSource(path /*isFile*/, true)
}

func (h *ReferenceResolutionHost) ReadFile(path string) (contents string, ok bool) {
	// Dont need to override as we cannot mimick read file
	return h.host.FS().ReadFile(path)
}

func (h *ReferenceResolutionHost) WriteFile(path string, data string, writeByteOrderMark bool) error {
	panic("should not be called by resolver")
}

// Removes `path` and all its contents. Will return the first error it encounters.
func (h *ReferenceResolutionHost) Remove(path string) error {
	panic("should not be called by resolver")
}

// DirectoryExists returns true if the path is a directory.
func (h *ReferenceResolutionHost) DirectoryExists(path string) bool {
	if h.host.FS().DirectoryExists(path) {
		h.handleDirectoryCouldBeSymlink(path)
		return true
	}
	return h.fileOrDirectoryExistsUsingSource(path /*isFile*/, false)
}

// GetAccessibleEntries returns the files/directories in the specified directory.
// If any entry is a symlink, it will be followed.
func (h *ReferenceResolutionHost) GetAccessibleEntries(path string) vfs.Entries {
	panic("should not be called by resolver")
}

func (h *ReferenceResolutionHost) Stat(path string) vfs.FileInfo {
	panic("should not be called by resolver")
}

// WalkDir walks the file tree rooted at root, calling walkFn for each file or directory in the tree.
// It is has the same behavior as [fs.WalkDir], but with paths as [string].
func (h *ReferenceResolutionHost) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	panic("should not be called by resolver")
}

// Realpath returns the "real path" of the specified path,
// following symlinks and correcting filename casing.
func (h *ReferenceResolutionHost) Realpath(path string) string {
	result, ok := h.symlinkCache.SymlinkedFiles()[h.toPath(path)]
	if ok {
		return result
	}
	return h.host.FS().Realpath(path)
}

func (h *ReferenceResolutionHost) toPath(path string) tspath.Path {
	return tspath.ToPath(path, h.GetCurrentDirectory(), h.UseCaseSensitiveFileNames())
}

func (h *ReferenceResolutionHost) handleDirectoryCouldBeSymlink(directory string) {
	if tspath.ContainsIgnoredPath(directory) {
		return
	}

	// Because we already watch node_modules, handle symlinks in there
	if !strings.Contains(directory, "/node_modules/") {
		return
	}
	directoryPath := tspath.Path(tspath.EnsureTrailingDirectorySeparator(string(h.toPath(directory))))
	if _, ok := h.symlinkCache.SymlinkedDirectories()[directoryPath]; ok {
		return
	}

	real := h.host.FS().Realpath(directory)
	var realPath tspath.Path
	if real == directory {
		// not symlinked
		h.symlinkCache.SetSymlinkedDirectory(directory, directoryPath, nil)
		return
	}
	if realPath = tspath.Path(tspath.EnsureTrailingDirectorySeparator(string(h.toPath(real)))); realPath == directoryPath {
		// not symlinked
		h.symlinkCache.SetSymlinkedDirectory(directory, directoryPath, nil)
		return
	}

	h.symlinkCache.SetSymlinkedDirectory(directory, directoryPath, &module.SymlinkedDirectory{
		Real:     tspath.EnsureTrailingDirectorySeparator(real),
		RealPath: realPath,
	})
}

func (h *ReferenceResolutionHost) fileOrDirectoryExistsUsingSource(fileOrDirectory string, isFile bool) bool {
	fileOrDirectoryExistsUsingSource := core.IfElse(isFile, h.fileExistsIfProjectReferenceDts, h.directoryExistsIfProjectReferenceDeclDir)
	// Check current directory or file
	result := fileOrDirectoryExistsUsingSource(fileOrDirectory)
	if result != core.TSUnknown {
		return result == core.TSTrue
	}

	// !!! sheetal this needs to be thread safe !!
	symlinkedDirectories := h.symlinkCache.SymlinkedDirectories()
	if symlinkedDirectories == nil {
		return false
	}
	fileOrDirectoryPath := tspath.ToPath(fileOrDirectory, h.GetCurrentDirectory(), h.UseCaseSensitiveFileNames())
	if !strings.Contains(string(fileOrDirectoryPath), "/node_modules/") {
		return false
	}
	if isFile {
		_, ok := h.symlinkCache.SymlinkedFiles()[fileOrDirectoryPath]
		if ok {
			return true
		}
	}

	// If it contains node_modules check if its one of the symlinked path we know of
	for directoryPath, symlinkedDirectory := range symlinkedDirectories {
		if symlinkedDirectory == nil {
			continue
		}

		relative, hasPrefix := strings.CutPrefix(string(fileOrDirectoryPath), string(directoryPath))
		if !hasPrefix {
			continue
		}
		result := fileOrDirectoryExistsUsingSource(string(symlinkedDirectory.RealPath) + relative)
		if result.IsTrue() {
			if isFile {
				// Store the real path for the file
				absolutePath := tspath.GetNormalizedAbsolutePath(fileOrDirectory, h.GetCurrentDirectory())
				h.symlinkCache.SetSymlinkedFile(
					fileOrDirectoryPath,
					symlinkedDirectory.Real+absolutePath[len(directoryPath):],
				)
			}
			return true
		}
	}
	return false
}

func (h *ReferenceResolutionHost) fileExistsIfProjectReferenceDts(file string) core.Tristate {
	source, ok := h.outputDtsToSource[tspath.ToPath(file, h.GetCurrentDirectory(), h.UseCaseSensitiveFileNames())]
	if ok {
		if source.Source != "" {
			return core.IfElse(h.host.FS().FileExists(source.Source), core.TSTrue, core.TSFalse)
		}
		return core.TSTrue
	}
	return core.TSUnknown
}

func (h *ReferenceResolutionHost) directoryExistsIfProjectReferenceDeclDir(dir string) core.Tristate {
	dirPath := tspath.ToPath(dir, h.GetCurrentDirectory(), h.UseCaseSensitiveFileNames())
	dirPathWithTrailingDirectorySeparator := dirPath + "/"
	for declDirPath := range h.dtsDirectories.Keys() {
		if dirPath == declDirPath ||
			// Any parent directory of declaration dir
			strings.HasPrefix(string(declDirPath), string(dirPathWithTrailingDirectorySeparator)) ||
			// Any directory inside declaration dir
			strings.HasPrefix(string(dirPath), string(declDirPath)+"/") {
			return core.TSTrue
		}
	}
	return core.TSUnknown
}
