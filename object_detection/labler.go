package object_detection

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/anish-krishnan/Tidepod/object_detection/models"
)

func GetLabelsForFile(filename string) ([]string, error) {
	model := models.NewCoco()
	err := model.Load()
	if err != nil {
		fmt.Printf("Error loading model: %v", err)
		panic(err)
	}

	defer model.CloseSession()

	file, err := os.Open("photo_storage/saved/" + filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	outcome := model.Predict(fileBytes)

	labelSet := map[string]bool{}
	for _, detection := range outcome.Detections {
		labelSet[detection.Label] = true
	}

	labels := []string{}
	for label := range labelSet {
		labels = append(labels, label)
	}

	return labels, nil
}
