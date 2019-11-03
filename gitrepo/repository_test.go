package gitrepo

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestOpenRepository(t *testing.T) {
	t.Run("should return error when directory does not exist", func(t *testing.T) {
		dir, err := ioutil.TempDir(os.TempDir(), "")
		require.Nil(t, err)
		defer os.RemoveAll(dir)

		repo, err := OpenGitRepository(dir)
		require.Nil(t, repo)
		require.Equal(t, fmt.Errorf("error opening repository: Could not find repository from '%s'", dir), err)
	})
}
