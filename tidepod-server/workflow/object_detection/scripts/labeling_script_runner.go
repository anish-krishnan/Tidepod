package script

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetLabelsWithPythonScript runs the imageai resnet classifier to
// label the image
func GetLabelsWithPythonScript(fullFilePath string) ([]string, error) {
	app := "python3"
	arg0 := "./workflow/object_detection/scripts/object_detection.py"
	arg1 := fullFilePath

	cmd := exec.Command(app, arg0, arg1)
	out, err := cmd.Output()

	if err != nil {
		fmt.Printf("Error running python object detection script (%s %s %s): %v", app, arg0, arg1, err)
		panic(err)
	}

	output := strings.Split(string(out), "\n")

	var labels []string

	for _, line := range output {
		if strings.HasPrefix(line, "RESULT") && len(line) > 7 {
			labels = strings.Split(line[7:], ";")
			break
		}
	}

	return labels, nil
}
