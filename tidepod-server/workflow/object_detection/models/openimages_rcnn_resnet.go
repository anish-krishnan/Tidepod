package models

import (
	"fmt"
	"path"

	"github.com/anish-krishnan/Tidepod/tidepod-server/workflow/object_detection/responses"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// Rcnn contains tensorflow model and human-readable string labels
type Rcnn struct {
	model  *tf.SavedModel
	labels []string
}

const fastModelPath = "../faster_rcnn_openimages_v4_inception_resnet_v2_1"

// NewRcnn returns a Rcnn object
func NewRcnn() *Rcnn {
	return &Rcnn{}
}

// Load loads the ssd_mobilenet_v1_Rcnn_2018_01_28 SavedModel.
func (c *Rcnn) Load() error {
	model, err := tf.LoadSavedModel(path.Join(basepath, fastModelPath, "/"), []string{""}, nil)
	if err != nil {
		return fmt.Errorf("Error loading model: %v", err)
	}
	c.model = model

	c.labels, err = readLabels(path.Join(basepath, fastModelPath, "labels.txt"))
	if err != nil {
		return fmt.Errorf("Error loading labels fileL %v", err)
	}

	return nil
}

// Predict predicts labels for an input image
func (c *Rcnn) Predict(data []byte) *responses.ObjectDetectionResponse {
	tensor, _ := makeTensorFromBytes(data)

	output, err := c.model.Session.Run(
		map[tf.Output]*tf.Tensor{c.model.Graph.Operation("image_tensor").Output(0): tensor},
		[]tf.Output{
			c.model.Graph.Operation("detection_boxes").Output(0),
			c.model.Graph.Operation("detection_classes").Output(0),
			c.model.Graph.Operation("detection_scores").Output(0),
			c.model.Graph.Operation("num_detections").Output(0),
		},
		nil,
	)

	if err != nil {
		fmt.Printf("Error running the session: %v", err)
		return nil
	}

	outcome := responses.NewObjectDetectionResponse(output, c.labels)
	return outcome
}

// CloseSession closes a session
func (c *Rcnn) CloseSession() {
	c.model.Session.Close()
}
