# Script to reset state of system

import os

filelist = [
    f
    for f in os.listdir(
        "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/tidepod-server/photo_storage/saved/"
    )
    if not f.startswith(".")
]
for f in filelist:
    os.remove(
        os.path.join(
            "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/tidepod-server/photo_storage/saved/",
            f,
        )
    )


filelist = [
    f
    for f in os.listdir(
        "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/tidepod-server/photo_storage/thumbnails/"
    )
    if not f.startswith(".")
]
for f in filelist:
    os.remove(
        os.path.join(
            "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/tidepod-server/photo_storage/thumbnails/",
            f,
        )
    )

filelist = [
    f
    for f in os.listdir(
        "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/tidepod-server/photo_storage/boxes/"
    )
    if not f.startswith(".")
]
for f in filelist:
    os.remove(
        os.path.join(
            "/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/tidepod-server/photo_storage/boxes/",
            f,
        )
    )


os.remove("/Users/anishkrishnan/src/github.com/anish-krishnan/Tidepod/tidepod-server/test.db")