package entity

type File struct{
	Name string
	Path string
	Checksum string
	Properties []string
}

func NewFile (name, path)(*File, error){
	file := &File{
		Name: name,
		path: path,
	}

	return file, nil
}

func (f *File) AddChecksum(path string)error{
	sha265sum := ""

	file, err := os.Open(filepath)
	if err != nil{
		return ErrOpeningOfFileUnsuccesfull
	}
	defer file.Close()
}