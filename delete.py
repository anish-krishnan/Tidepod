# Script to reset state of system

import os

filelist = [
    f
    for f in os.listdir(
        "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/photo_storage/saved/"
    )
    if not f.startswith(".")
]
for f in filelist:
    os.remove(
        os.path.join(
            "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/photo_storage/saved/",
            f,
        )
    )


filelist = [
    f
    for f in os.listdir(
        "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/photo_storage/thumbnails/"
    )
    if not f.startswith(".")
]
for f in filelist:
    os.remove(
        os.path.join(
            "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/photo_storage/thumbnails/",
            f,
        )
    )


os.remove("/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/test.db")