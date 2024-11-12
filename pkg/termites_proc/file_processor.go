package termites_proc

import (
	"fmt"
	"github.com/RoelofRuis/termites/pkg/termites"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type FileProcessor[In any, Out any] struct {
	DataIn    *termites.InPort
	ResultOut *termites.OutPort

	procFunc func(In) (Out, error)

	rootDir           string
	cleanupExtensions []string
	cleanupFreq       time.Duration
	cleanupActive     bool
	cleanupFileCount  int
}

func NewFileProcessor[In any, Out any](
	procFunc func(In) (Out, error),
	rootDir string,
	cleanupExtensions []string,
	cleanupActive bool,
	cleanupFreq time.Duration,
	cleanupFileCount int,
) *FileProcessor[In, Out] {
	builder := termites.NewBuilder("FileProcessor")

	n := &FileProcessor[In, Out]{
		DataIn:            termites.NewInPort[In](builder),
		ResultOut:         termites.NewOutPort[Out](builder),
		procFunc:          procFunc,
		rootDir:           rootDir,
		cleanupExtensions: cleanupExtensions,
		cleanupFreq:       cleanupFreq,
		cleanupActive:     cleanupActive,
		cleanupFileCount:  cleanupFileCount,
	}

	builder.OnRun(n.Run)

	return n
}

func (f *FileProcessor[In, Out]) Run(e termites.NodeControl) error {
	err := f.cleanAll()
	if err != nil {
		return err
	}

	cleanup := time.NewTicker(f.cleanupFreq)
	defer cleanup.Stop()

	for {
		select {
		case <-cleanup.C:
			err := f.cleanOldest()
			if err != nil {
				return err
			}

		case msg := <-f.DataIn.Receive():
			data := msg.Data.(In)
			out, err := f.procFunc(data)
			if err != nil {
				e.LogError("processing error", err)
				continue
			}

			f.ResultOut.Send(out)
		}
	}
}

func (f *FileProcessor[In, Out]) cleanAll() error {
	files, err := os.ReadDir(f.rootDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && f.isRemovable(file.Name()) {
			err := f.remove(file.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *FileProcessor[In, Out]) cleanOldest() error {
	files, err := os.ReadDir(f.rootDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}
	if len(files) == f.cleanupFileCount {
		fileInfos := make([]os.FileInfo, 0, len(files))
		for _, file := range files {
			if file.IsDir() || !f.isRemovable(file.Name()) {
				continue
			}

			info, err := file.Info()
			if err != nil {
				fmt.Printf("error reading file info: %e\n", err)
				continue
			}
			fileInfos = append(fileInfos, info)
		}

		sort.Slice(fileInfos, func(i, j int) bool {
			return fileInfos[i].ModTime().Before(fileInfos[j].ModTime())
		})

		excess := len(files) - f.cleanupFileCount

		for i := 0; i < excess; i++ {
			file := fileInfos[i]
			if err := f.remove(file.Name()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *FileProcessor[In, Out]) remove(filename string) error {
	path := filepath.Join(f.rootDir, filename)
	if !f.cleanupActive {
		fmt.Printf("Would remove %s\n", path)
		return nil
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove %s: %w", path, err)
	}

	return nil
}

func (f *FileProcessor[In, Out]) isRemovable(filename string) bool {
	for _, extension := range f.cleanupExtensions {
		if strings.HasSuffix(filename, extension) {
			return true
		}
	}
	return false
}
