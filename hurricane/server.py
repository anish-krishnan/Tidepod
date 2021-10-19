import tornado.ioloop
import tornado.web
import json
from imageai.Detection import ObjectDetection
from mtcnn import MTCNN

model_path = "./resnet50_coco_best_v2.1.0.h5"
detector = ObjectDetection()
detector.setModelTypeAsRetinaNet()
detector.setModelPath(model_path)
detector.loadModel()

from matplotlib import pyplot

faceDetector = MTCNN()


class LabelImageHandler(tornado.web.RequestHandler):
    def post(self):
        filename = self.get_argument("filename")
        labels = self.getLabelsForImage(filename)
        print({k: self.get_argument(k) for k in self.request.arguments}, labels)
        self.write(json.dumps({"labels": labels}))

    def getLabelsForImage(self, filename):
        input_path = "../tidepod-server/photo_storage/TEMP/" + filename
        output_path = "../tidepod-server/photo_storage/TEMP/out.jpg"

        detection = detector.detectObjectsFromImage(
            input_image=input_path, output_image_path=output_path
        )
        THRESHOLD = 60
        labels = set()
        for item in detection:
            label, confidence = item["name"], item["percentage_probability"]
            if confidence > THRESHOLD:
                labels.add(label)
        return list(labels)


class FaceDetectHandler(tornado.web.RequestHandler):
    def post(self):
        filename = self.get_argument("filename")
        face_locations = self.getFaceLocationsForImage2(filename)
        print(filename, face_locations)
        self.write(json.dumps({"faces": face_locations}))

    def getFaceLocationsForImage2(self, filename):
        input_path = "../tidepod-server/photo_storage/TEMP/" + filename

        pixels = pyplot.imread(input_path)
        faces = faceDetector.detect_faces(pixels)

        # filter on confidence
        faces = [face for face in faces if face["confidence"] >= 0.6]

        face_locations = [face["box"] for face in faces]

        face_locations = [
            {"minX": x, "minY": y, "maxX": x + w, "maxY": y + h}
            for x, y, w, h in face_locations
        ]

        return face_locations

    def getFaceLocationsForImage(self, filename):
        input_path = "../tidepod-server/photo_storage/TEMP/" + filename

        image = face_recognition.load_image_file(input_path)
        face_locations = face_recognition.face_locations(
            image, number_of_times_to_upsample=0, model="cnn"
        )

        faces = [
            {"minX": left, "minY": top, "maxX": right, "maxY": bottom}
            for top, right, bottom, left in face_locations
        ]

        return faces


def make_app():
    return tornado.web.Application(
        [
            (r"/labelImage", LabelImageHandler),
            (r"/faceDetect", FaceDetectHandler),
        ]
    )


if __name__ == "__main__":
    app = make_app()
    app.listen(2999)
    tornado.ioloop.IOLoop.current().start()