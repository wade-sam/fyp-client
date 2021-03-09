package backup

import "github.com/wade-sam/fypclient/entity"

type BSRepository interface {
	DirectoryFile(Directories *entity.Directory)
	ExpectedFiles(Directories *entity.Directory) error
	FileBackupMesssage(file *entity.File) error
}

type WFRepository interface {
	GetbackupFile() (*entity.Directory, error)
	CreateBackupFile(*entity.Directory) error
}

type SNRepository interface {
	BackupLayout(Directories *entity.Directory) error
	//ExpectedFiles(Directories *entity.Directory)error //Backupserver will send a file report to the storage node
	FileBackup(file []byte) error
}

// type Repository interface {
// 	BSRepository
// 	WFRepository
// 	SNRepository
// }

type Usecase interface {
	StartIncrementalBackup(ignore []string) error
	StartFullBackup(ignore []string) error
}
