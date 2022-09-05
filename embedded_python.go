package go_embed_python

import (
	"fmt"
	"github.com/kluctl/go-embed-python/embed_util"
	"github.com/kluctl/go-embed-python/internal/embed"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type EmbeddedPython struct {
	e *embed_util.EmbeddedFiles
}

// NewEmbeddedPython creates a new EmbeddedPython instance. The embedded source code and python binaries are
// extracted on demand using the given name as the base for the temporary directory.
func NewEmbeddedPython(name string) (*EmbeddedPython, error) {
	e, err := embed_util.NewEmbeddedFiles(embed.PythonLib, fmt.Sprintf("python-%s", name))
	if err != nil {
		return nil, err
	}
	return &EmbeddedPython{
		e: e,
	}, nil
}

func NewEmbeddedPythonWithTmpDir(tmpDir string) (*EmbeddedPython, error) {
	e, err := embed_util.NewEmbeddedFilesWithTmpDir(embed.PythonLib, tmpDir)
	if err != nil {
		return nil, err
	}
	return &EmbeddedPython{
		e: e,
	}, nil
}

func (ep *EmbeddedPython) Cleanup() error {
	return ep.e.Cleanup()
}

func (ep *EmbeddedPython) GetExtractedPath() string {
	return ep.e.GetExtractedPath()
}

func (ep *EmbeddedPython) GetBinPath() string {
	if runtime.GOOS == "windows" {
		return ep.GetExtractedPath()
	} else {
		return filepath.Join(ep.GetExtractedPath(), "bin")
	}
}

func (ep *EmbeddedPython) GetExePath() string {
	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	} else {
		suffix = "3"
	}
	return filepath.Join(ep.GetBinPath(), "python"+suffix)
}

func (ep *EmbeddedPython) PythonCmd(args ...string) *exec.Cmd {
	return ep.PythonCmd2(args)
}

func (ep *EmbeddedPython) PythonCmd2(args []string) *exec.Cmd {
	exePath := ep.GetExePath()

	cmd := exec.Command(exePath, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PYTHONHOME=%s", ep.GetExtractedPath()))

	return cmd
}
