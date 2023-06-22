package publish

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

// Return a list of files inside the given directory, ignoring the given globs
func (cmd *PublishCmd) getFileList(root string, ignore_globs []string) ([]string, error) {
	var files []string

	// parse ignore globs with gitignore.Pattern.ParsePattern
	var patterns []gitignore.Pattern
	for _, glob := range ignore_globs {
		pattern := gitignore.ParsePattern(glob, []string{})
		patterns = append(patterns, pattern)
	}
	matcher := gitignore.NewMatcher(patterns)

	err := cmd.FilepathWalk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// use relative path to compare to ignore globs
		rel_path, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		rel_path_list := strings.Split(rel_path, string(filepath.Separator))
		ignore_matched := matcher.Match(rel_path_list, info.IsDir())
		if !ignore_matched {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
