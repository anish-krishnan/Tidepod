import tornado.ioloop
import tornado.web
import json
from imageai.Detection import ObjectDetection

model_path = "./resnet50_coco_best_v2.1.0.h5"
detector = ObjectDetection()
detector.setModelTypeAsRetinaNet()
detector.setModelPath(model_path)
detector.loadModel()


class MainHandler(tornado.web.RequestHandler):
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


def make_app():
    return tornado.web.Application(
        [
            (r"/", MainHandler),
        ]
    )


if __name__ == "__main__":
    app = make_app()
    app.listen(2999)
    tornado.ioloop.IOLoop.current().start()