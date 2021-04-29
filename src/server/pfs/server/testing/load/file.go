package load

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"path"

	"github.com/pachyderm/pachyderm/v2/src/internal/storage/chunk"
	"github.com/pachyderm/pachyderm/v2/src/internal/uuid"
)

// MemFile contains file data and metadata. It's generated by FileSource.Next(),
// and used for fuzzer-generated PutFile operations
type MemFile struct {
	path    string
	content []byte
}

// NewMemFile constructs a MemFile at the given 'path' containing 'data'.
func NewMemFile(path string, data []byte) *MemFile {
	return &MemFile{
		path:    path,
		content: data,
	}
}

// Path returns the path at which a given MemFile will be written.
func (mf *MemFile) Path() string {
	return mf.path
}

// Reader returns an io.Reader that produces 'mf's contents (which are consumed
// by PutFile).
func (mf *MemFile) Reader() io.Reader {
	return bytes.NewReader(mf.content)
}

// FileSourceSpec configures a FileSource. FileSourceSpecs are included in the
// CommitsSpec passed to NewEnv during fuzzer initialization; they're used to
// set up the env's FileSources, which are in turn used throughout the fuzz test
// by PutFile operations.
type FileSourceSpec struct {
	// Name is a user-provided name for this FileSource, e.g. "random",
	// "testdataset", etc.
	Name string `yaml:"name,omitempty"`

	// RandomFileSourceSpec, if set, configures the resulting FileSource to be a
	// Random-type FileSource (files consist of randomly-generated content).
	RandomFileSourceSpec *RandomFileSourceSpec `yaml:"random,omitempty"`
}

// FileSource is an interface for file-generators. FileSources are part of the
// load-testing environment and are used to provide Files to fuzzer-generated
// PutFile operations.
type FileSource interface {
	// Next generates and returns a new file.
	Next() *MemFile
}

// NewFileSource initializes a FileSource from a FileSourceSpec (called during
// intialization, this sets up the load-testing environment's FileSources).
func NewFileSource(spec *FileSourceSpec) FileSource {
	return newRandomFileSource(spec.RandomFileSourceSpec)
}

// RandomFileSourceSpec specifies a RandomFile-type FileSource
type RandomFileSourceSpec struct {
	// IncrementPath, if true, places each generated file underneath a directory
	// whose name increments on each call to Next();
	IncrementPath bool `yaml:"incrementPath,omitempty"`

	// RandomDirectorySpec, if set, places each generated file (and each
	// incrementally-named directory, if IncrementPath is set) underneath a
	// directory named with a UUID.
	RandomDirectorySpec *RandomDirectorySpec `yaml:"directory,omitempty"`

	// FuzzSizeSpecs specify the possible sizes of any file generated using this
	// RandomFileSourceSpec. When generating a file, a FuzzSizeSpec is selected
	// from this list with probably equal to its 'Prob' field; the 'Prob' fields
	// of all FuzzSizeSpecs here must sum to 1.
	FuzzSizeSpecs []*FuzzSizeSpec `yaml:"fuzzSize,omitempty"`
}

type randomFileSource struct {
	spec      *RandomFileSourceSpec
	dirSource *randomDirectorySource
	next      int64
}

func newRandomFileSource(spec *RandomFileSourceSpec) FileSource {
	var dirSource *randomDirectorySource
	if spec.RandomDirectorySpec != nil {
		dirSource = &randomDirectorySource{
			spec: spec.RandomDirectorySpec,
		}
	}
	return &randomFileSource{
		spec:      spec,
		dirSource: dirSource,
	}
}

func (rfs *randomFileSource) Next() *MemFile {
	sizeSpec := FuzzSize(rfs.spec.FuzzSizeSpecs)
	min, max := sizeSpec.Min, sizeSpec.Max
	size := min
	if max > min {
		size += rand.Intn(max - min)
	}
	return NewMemFile(rfs.nextPath(), chunk.RandSeq(size))
}

func (rfs *randomFileSource) nextPath() string {
	var dir string
	if rfs.dirSource != nil {
		dir = rfs.dirSource.nextPath()
	}
	if rfs.spec.IncrementPath {
		next := rfs.next
		rfs.next += 1
		return path.Join(dir, fmt.Sprintf("%016d", next))
	}
	return path.Join(dir, uuid.NewWithoutDashes())
}

type RandomDirectorySpec struct {
	Depth int   `yaml:"depth,omitempty"`
	Run   int64 `yaml:"run,omitempty"`
}

type randomDirectorySource struct {
	spec *RandomDirectorySpec
	next string
	run  int64
}

func (rds *randomDirectorySource) nextPath() string {
	if rds.next == "" {
		depth := rand.Intn(rds.spec.Depth)
		for i := 0; i < depth; i++ {
			rds.next = path.Join(rds.next, uuid.NewWithoutDashes())
		}
	}
	dir := rds.next
	rds.run++
	if rds.run == rds.spec.Run {
		rds.next = ""
	}
	return dir

}

// FilesSpec specifi
type FilesSpec struct {
	// Count specifies the number of files that Files will write during a PutFile
	// operation
	Count int `yaml:"count,omitempty"`

	// FuzzFileSpecs specify the possible files that may be generated using this
	// FilesSpec. When generating a file, a FuzzFileSpec is selected from this
	// list with probably equal to its 'Prob' field; the 'Prob' fields of all
	// FuzzFileSpecs here must sum to 1.
	FuzzFileSpecs []*FuzzFileSpec `yaml:"fuzzFile,omitempty"`
}

func Files(env *Env, spec *FilesSpec) ([]*MemFile, error) {
	var files []*MemFile
	for i := 0; i < spec.Count; i++ {
		file, err := FuzzFile(env, spec.FuzzFileSpecs)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

// TODO: Add different types of files.
type FileSpec struct {
	Source string `yaml:"source,omitempty"`
}

// File generates and returns a file
func File(env *Env, spec *FileSpec) (*MemFile, error) {
	return env.FileSource(spec.Source).Next(), nil
}

// SizeSpec specifies the possible size range of a PutFileRequest's Value field
// (in bytes)
type SizeSpec struct {
	Min int `yaml:"min,omitempty"`
	Max int `yaml:"max,omitempty"`
}
