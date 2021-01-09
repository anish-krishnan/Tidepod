from imageai.Detection import ObjectDetection
import sys


input_path = sys.argv[1]
model_path = "./workflow/object_detection/scripts/resnet50_coco_best_v2.1.0.h5"
output_path = "./photo_storage/TEMP/out.jpg"

detector = ObjectDetection()
detector.setModelTypeAsRetinaNet()
detector.setModelPath(model_path)

detector.loadModel()

detection = detector.detectObjectsFromImage(input_image=input_path, output_image_path=output_path)

THRESHOLD = 60
labels = set()
for item in detection:
  label, confidence = item["name"], item["percentage_probability"]
  if confidence > THRESHOLD:
    labels.add(label)

result = "RESULT " + ";".join(labels)
print(result)


