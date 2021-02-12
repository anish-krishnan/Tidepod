# Script to reset state of system

import os

# Remove all saved images
directories = [
    "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/tidepod-server/photo_storage/saved/",
    "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/tidepod-server/photo_storage/thumbnails/",
    "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/tidepod-server/photo_storage/boxes/",
    "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/tidepod-server/photo_storage/TEMP/",
]

for directory in directories:
    filelist = [file for file in os.listdir(directory) if not file.startswith(".")]
    for file in filelist:
        os.remove(os.path.join(directory, file))

# Remove database
os.remove("/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/tidepod-server/test.db")