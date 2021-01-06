package workflow

import (
	"fmt"
	"os/exec"

	"github.com/anish-krishnan/Tidepod/tidepod-server/entity"
	"github.com/disintegration/imaging"
	"gorm.io/gorm"
)

// CreateVideoThumbnail takes a video, and creates a 200x200 thumbnail
// and saves it to the photo_storage/thumbnails/ directory.
// It also rotates the thumbnail as needed
func CreateVideoThumbnail(db *gorm.DB, photo *entity.Photo) {
	thumbFilename := fmt.Sprintf("%d%s", photo.ID, ".jpg")

	photo.ThumbnailFilePath = thumbFilename
	db.Save(photo)

	// Save first frame from video
	app := "ffmpeg"
	arg0 := "-i"
	arg1 := "photo_storage/saved/" + photo.FilePath
	arg2 := "-ss"
	arg3 := "00:00:00.001"
	arg4 := "-vframes"
	arg5 := "1"
	arg6 := "photo_storage/thumbnails/" + thumbFilename

	cmd := exec.Command(app, arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	_, err := cmd.Output()

	if err != nil {
		panic(err)
	}

	// Convert first frame to 200x200 thumbnail
	img, err := imaging.Open("photo_storage/thumbnails/" + thumbFilename)

	if err != nil {
		panic(err)
	}
	thumb := imaging.Thumbnail(img, 200, 200, imaging.CatmullRom)

	err = imaging.Save(thumb, "photo_storage/thumbnails/"+thumbFilename)
	if err != nil {
		fmt.Println("ERROR saving thumbnail", "photo_storage/thumbnails/"+thumbFilename)
		panic(err)
	}
}
