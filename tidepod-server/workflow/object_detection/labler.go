package objectdetection

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/anish-krishnan/Tidepod/tidepod-server/workflow/object_detection/models"
)

// GetLabelsForFile converts a txt file containing a label on each new line
// to a string array
func GetLabelsForFile(fullFilePath string) ([]string, error) {
	model := models.NewCoco()
	err := model.Load()
	if err != nil {
		fmt.Printf("Error loading model: %v", err)
		panic(err)
	}

	defer model.CloseSession()

	file, err := os.Open(fullFilePath)
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
