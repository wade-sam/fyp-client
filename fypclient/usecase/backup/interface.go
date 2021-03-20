package backup

import "github.com/wade-sam/fypclient/entity"

type BSRepository interface {
	DirectoryFile(Directories *entity.Directory)
	ExpectedFiles(map[string]*entity.File) error
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
	BackupDirectoryScan(head string, ignore map[string]string) (*map[string]interface{}, error)
	FullBackupCopy(n *entity.Directory, chn chan (FileTransfer))
	//StartIncrementalBackup(ignore []string) error
	//StartFullBackup(ignore []string) error
}
