package afero

import (
	"os"
	"time"
)

// An afero Fs with an extra filter
//
// The FilterFs is run before the source Fs, any non nil error is returned
// to the caller without going to the source Fs. If every filter in the
// chain returns a nil error, the call is sent to the source Fs.
//
// see the TestReadonlyRemoveAndRead() in filter_test.go for an example use
// of filtering (e.g. admins get write access, normal users just readonly)
type FilterFs interface {
	Fs
	AddFilter(Fs)
}

type Filter struct {
	chain  []Fs
	source Fs
}

// create a new FilterFs that implements Fs, argument must be an Fs, not
// a FilterFs
func NewFilter(fs Fs) FilterFs {
	return &Filter{source: fs}
}

// prepend a filter in the filter chain
func (f *Filter) AddFilter(fs Fs) {
	c := []Fs{fs}
	for _, ch := range f.chain {
		c = append(c, ch)
	}
	f.chain = c
}

func (f *Filter) Create(name string) (file File, err error) {
	for _, c := range f.chain {
		file, err = c.Create(name)
		if err != nil {
			return
		}
	}
	return f.source.Create(name)
}

func (f *Filter) Mkdir(name string, perm os.FileMode) (err error) {
	for _, c := range f.chain {
		err = c.Mkdir(name, perm)
		if err != nil {
			return
		}
	}
	return f.source.Mkdir(name, perm)
}

func (f *Filter) MkdirAll(path string, perm os.FileMode) (err error) {
	for _, c := range f.chain {
		err = c.MkdirAll(path, perm)
		if err != nil {
			return
		}
	}
	return f.source.MkdirAll(path, perm)
}

func (f *Filter) Open(name string) (file File, err error) {
	for _, c := range f.chain {
		file, err = c.Open(name)
		if err != nil {
			return
		}
	}
	return f.source.Open(name)
}

func (f *Filter) OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	for _, c := range f.chain {
		file, err = c.OpenFile(name, flag, perm)
		if err != nil {
			return
		}
	}
	return f.source.OpenFile(name, flag, perm)
}

func (f *Filter) Remove(name string) (err error) {
	for _, c := range f.chain {
		err = c.Remove(name)
		if err != nil {
			return
		}
	}
	return f.source.Remove(name)
}

func (f *Filter) RemoveAll(path string) (err error) {
	for _, c := range f.chain {
		err = c.RemoveAll(path)
		if err != nil {
			return
		}
	}
	return f.source.RemoveAll(path)
}

func (f *Filter) Rename(oldname, newname string) (err error) {
	for _, c := range f.chain {
		err = c.Rename(oldname, newname)
		if err != nil {
			return
		}
	}
	return f.source.Rename(oldname, newname)
}

func (f *Filter) Stat(name string) (fi os.FileInfo, err error) {
	for _, c := range f.chain {
		fi, err = c.Stat(name)
		if err != nil {
			return
		}
	}
	return f.source.Stat(name)
}

func (f *Filter) Name() string {
	return f.source.Name()
}

func (f *Filter) Chmod(name string, mode os.FileMode) (err error) {
	for _, c := range f.chain {
		err = c.Chmod(name, mode)
		if err != nil {
			return
		}
	}
	return f.source.Chmod(name, mode)
}

func (f *Filter) Chtimes(name string, atime, mtime time.Time) (err error) {
	for _, c := range f.chain {
		err = c.Chtimes(name, atime, mtime)
		if err != nil {
			return
		}
	}
	return f.source.Chtimes(name, atime, mtime)
}
