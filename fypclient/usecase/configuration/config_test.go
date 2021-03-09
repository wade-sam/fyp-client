package configuration_test

import (
	"compress/gzip"
	//S"encoding/json"
	"fmt"

	jsoniter "github.com/json-iterator/go"

	//"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	repo "github.com/wade-sam/fypclient/Infrastructure/Repositories/writetofile"
	"github.com/wade-sam/fypclient/entity"
	service "github.com/wade-sam/fypclient/usecase/configuration"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Test_DirectoryScan(t *testing.T) {
	cfgrepo := repo.NewFileRepo()
	newservice := service.NewConfigurationService(cfgrepo)
	ignore := make(map[string]string)
	ignore["/proc"] = "/proc"
	ignore["/bin"] = "/bin"
	directory, err := newservice.DirectoryScan("/")
	assert.Nil(t, err)
	//fmt.Println(directory)

	outputdirectory, err := json.MarshalIndent(directory, "", "   ")
	if err != nil {
		fmt.Println(entity.ErrCouldNotMarshallJSON)
	}
	di, err := os.Create("directory.gzip")
	q := gzip.NewWriter(di)
	q.Write([]byte(outputdirectory))
	q.Close()
	//err = ioutil.WriteFile("directory.json", outputfile, 0775)
	//	if err != nil {
	//		fmt.Println(entity.ErrCouldNotWriteToFile)
	//	}

}

// func Test_GetDetails(t *testing.T) {
// 	cfgrepo := repo.NewFileRepo()
// 	newservice := service.NewConfigurationService(cfgrepo)
// 	client, err := newservice.GetClientName()
// 	assert.Nil(t, err)
// 	assert.Equal(t, "george", client)

// }

// func Test_SetClientName(t *testing.T) {
// 	cfgrepo := repo.NewFileRepo()
// 	newservice := service.NewConfigurationService(cfgrepo)
// 	client, err := newservice.GetStorageNode()
// 	assert.Nil(t, err)
// 	err = newservice.SetStorageNode("jack")
// 	assert.Nil(t, err)
// 	client1, err := newservice.GetClientName()
// 	assert.NotEqual(t, client, client1)
// }
