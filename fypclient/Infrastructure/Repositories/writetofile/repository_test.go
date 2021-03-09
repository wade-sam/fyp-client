package writetofile_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	repo "github.com/wade-sam/fypclient/Infrastructure/Repositories/writetofile"

	//repo "github.com/wade-sam/fypclient/Infrastructure/Repositiories/writetofile"
	"github.com/wade-sam/fypclient/entity"
)

func structsInitialise() *repo.FileStruct {
	rabbit := repo.RabbitConfig{
		Schema:         "test",
		Username:       "user",
		Password:       "pass",
		Host:           "192.168.1.1",
		Port:           "8976",
		VHost:          "/",
		ConnectionName: "wadee",
	}

	file := repo.FileStruct{
		ClientName:    "sam",
		BackupServer:  "192.168.1.1",
		StorageNode:   "192.168.1.2",
		RabbitDetails: &rabbit,
	}
	return &file
}
func initialiseFile(file *repo.FileStruct) error {
	outputFile, err := json.MarshalIndent(file, "", "	")
	if err != nil {
		return entity.ErrCouldNotMarshallJSON
	}
	err = ioutil.WriteFile("config.json", outputFile, 0775)
	if err != nil {
		return entity.ErrCouldNotWriteToFile
	}
	return nil
}

// func Test_InitialiseRepo(t *testing.T) {
// 	//fr := repo.NewFileRepo()
// 	file := structsInitialise()
// 	err := initialiseFile(file)
// 	assert.Nil(t, err)
// }

func Test_GetClientName(t *testing.T) {
	fr := repo.NewFileRepo()
	//prefile := structsInitialise()
	_, err := fr.GetClientName()
	assert.Nil(t, err)
	//assert.Equal(t, prefile.ClientName, name)
}

func Test_WriteClientName(t *testing.T) {
	fr := repo.NewFileRepo()
	name, err := fr.GetClientName()
	assert.Nil(t, err)
	err = fr.SetClientName("george")
	assert.Nil(t, err)
	name2, err := fr.GetClientName()
	assert.Nil(t, err)
	assert.Equal(t, name, name2)
	fmt.Println(name, name2)
}

func Test_GetRabbitDetails(t *testing.T) {
	fr := repo.NewFileRepo()
	file := structsInitialise()
	rabbit, err := fr.GetRabbitDetails()
	assert.Nil(t, err)
	assert.Equal(t, file.RabbitDetails, rabbit)
	fmt.Println(rabbit)
}
