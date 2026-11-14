package filesystemcap

import "github.com/gzjjjfree/v5-result/common/platform/filesystem/fsifce"

type FileSystemCapabilitySet interface {
	OpenFileForReadSeek() fsifce.FileSeekerFunc
	OpenFileForRead() fsifce.FileReaderFunc
	OpenFileForWrite() fsifce.FileWriterFunc
	ReadDir() fsifce.FileReadDirFunc
	RemoveFile() fsifce.FileRemoveFunc
}
