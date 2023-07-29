package upbin

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
)

func TestUpdateBin(t *testing.T) {
	notifyFunc = func() (bool, error) {
		return true, nil
	}

	needToUpdateFunc = func() bool {
		return true
	}

	wantToUpdateFunc = func() bool {
		return true
	}

	prepareURLFunc = func() (string, error) {
		return "http://example.com/file", nil
	}

	managersPackagesFunc = func() (bool, error) {
		return false, nil
	}

	whichFunc = func(string) (string, error) {
		return "/path/to/azioncli", nil
	}

	downloadFileFunc = func(string) error {
		return nil
	}

	replaceFileFunc = func(string) error {
		return nil
	}

	err := UpdateBin()
	assert.NoError(t, err)
}

type MockReferenceIter struct {
	References []*plumbing.Reference
}

func (m *MockReferenceIter) ForEach(fn func(*plumbing.Reference) error) error {
	for _, ref := range m.References {
		err := fn(ref)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestLatestTag(t *testing.T) {
	// Create mock references
	refs := []*plumbing.Reference{
		plumbing.NewReferenceFromStrings("refs/tags/0.1.0", "commit1"),
		plumbing.NewReferenceFromStrings("refs/tags/0.2.0", "commit2"),
		plumbing.NewReferenceFromStrings("refs/tags/0.3.0", "commit3"),
		plumbing.NewReferenceFromStrings("refs/tags/0.3.1", "commit4"),
		plumbing.NewReferenceFromStrings("refs/tags/0.4.0", "commit5"),
		plumbing.NewReferenceFromStrings("refs/tags/0.4.1", "commit6"),
		plumbing.NewReferenceFromStrings("refs/tags/0.4.3", "commit7"),
		plumbing.NewReferenceFromStrings("refs/tags/1.0.0", "commit8"),
		plumbing.NewReferenceFromStrings("refs/tags/1.1.0", "commit9"),
		plumbing.NewReferenceFromStrings("refs/tags/1.1.1", "commit10"),
	}

	// Create mock reference iterator
	mockIter := &MockReferenceIter{
		References: refs,
	}

	// Call the function under test
	tag, err := latestTag(mockIter)

	// Assert the expected result
	assert.NoError(t, err)
	assert.Equal(t, "refs/tags/1.1.1", tag)
}

func TestGetLastActivity(t *testing.T) {
	openConfigFunc = func() (map[string]interface{}, error) {
		return map[string]interface{}{
			"LAST_ACTIVITY": "2021-10-01",
		}, nil
	}

	expectedLastActivity := "2021-10-01"

	lastActivity, err := getLastActivity()
	assert.NoError(t, err)
	assert.Equal(t, expectedLastActivity, lastActivity)
}
