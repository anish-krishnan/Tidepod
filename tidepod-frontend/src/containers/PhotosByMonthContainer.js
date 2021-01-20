import React, { useCallback } from 'react'
import Gallery from 'react-photo-gallery';

import { FaPlayCircle } from 'react-icons/fa';
import ProgressBar from 'react-bootstrap/ProgressBar'
import InfiniteScroll from "react-infinite-scroll-component";
import placeholder from '../placeholder.jpg'
import {
  ReasonPhrases,
  StatusCodes,
  getReasonPhrase,
  getStatusCode,
} from 'http-status-codes';


class PhotosByMonthContainer extends React.Component {
  state = {
    photosByMonth: [],
    files: [],

    uploading: false,
    currentUploadNum: 0,
    uploadLength: 0,
    uploadErrors: [],

    numPhotos: 0,
    offset: 0,
    hasMorePhotos: true
  }

  componentDidMount() {
    fetch("/api/photosByMonth/" + this.state.offset, {
      headers: { "Token": this.props.idToken }
    })
      .then(resp => resp.json())
      .then(photosByMonth => {
        this.setState({
          photosByMonth: photosByMonth,
          offset: this.state.offset + 1,
          numPhotos: photosByMonth.map(p => p.Photos ? p.Photos.length : 0).reduce((acc, val) => acc + val, 0),
        })
      })
  }

  onChange = (event) => {
    this.setState({
      files: event.target.files
    })
  }

  handleUpload = async (event) => {
    for (var i = 0; i < this.state.files.length; i++) {
      this.setState({
        uploading: true,
        currentUploadNum: i + 1,
        uploadLength: this.state.files.length
      })

      var file = this.state.files[i]
      const formData = new FormData();
      formData.append('files', file)
      formData.append('infoArray', file.lastModified)
      formData.append('token', this.props.idToken)

      await fetch("/api/upload", {
        method: 'POST',
        body: formData,
      })
        .then(resp => {
          if (resp.status == StatusCodes.UNSUPPORTED_MEDIA_TYPE) {
            this.setState({
              uploadErrors: this.state.uploadErrors.concat([{ "filename": file.name, "error": ReasonPhrases.UNSUPPORTED_MEDIA_TYPE }])
            })
          }
          return resp.json()
        })
        .then(response => {
          console.log("response", response)
        })
        .catch(err => {
          console.log("error", err)
        })
    }

    this.setState({
      uploading: false,
    })
  }

  onClick = (event) => {
    const { history } = this.props;
    if (history) {
      history.push({
        pathname: '/photo/' + event.target.id,
        state: { idToken: this.props.idToken }
      });
    }
  }

  imageRenderer = ({ index, left, top, key, photo }) => {
    const [height, width] = [photo.height, photo.width]
    if (photo.photo.MediaType == "photo") {
      return (

        <div style={{ display: "inline" }}>
          <img
            onClick={this.onClick}
            id={photo.id}
            src={key}
            style={{ margin: "2px", height: { height }, width: { width }, display: "block", cursor: "pointers" }} />
        </div>

      )
    } else {
      return (
        <span style={{ margin: "2px", position: "relative", height: { height }, width: { width }, display: "inline", cursor: "pointers" }}>
          <img
            onClick={this.onClick}
            id={photo.id}
            src={key}
          />
          <h1><FaPlayCircle style={{ position: "absolute", top: "0%", left: "80%", color: "white" }} /></h1>
        </span >
      )
    }
  }

  fetchMoreData = () => {
    console.log("fetchMoreData", this.state.offset)

    fetch("/api/photosByMonth/" + this.state.offset, {
      headers: { "Token": this.props.idToken }
    })
      .then(resp => resp.json())
      .then(newPhotos => {

        let original = this.state.photosByMonth
        let numNewPhotos = newPhotos.map(p => p.Photos ? p.Photos.length : 0).reduce((acc, val) => acc + val, 0)
        // debugger;

        console.log("new photos", newPhotos)
        if (newPhotos.length == 1 && newPhotos.[0].Photos == null) {
          this.setState({
            hasMorePhotos: false
          })
          console.log("DONE!")
          return
        }


        if (original[original.length - 1].Month === newPhotos[0].Month) {
          original[original.length - 1].Photos = original[original.length - 1].Photos.concat(newPhotos[0].Photos)
          this.setState({
            photosByMonth: original.concat(newPhotos.slice(1)),
            offset: this.state.offset + 1,
            numPhotos: this.state.numPhotos + numNewPhotos
          })
        } else {
          this.setState({
            photosByMonth: original.concat(newPhotos),
            offset: this.state.offset + 1,
            numPhotos: this.state.numPhotos + numNewPhotos
          })
        }
      })

    console.log("returning")

  }


  render() {
    const uploadPercent = this.state.uploadLength ? (this.state.currentUploadNum / this.state.uploadLength * 100) : 0
    const uploadBarLabel = (this.state.currentUploadNum) + " of " + (this.state.uploadLength)

    return (

      <div className="main-container" >
        {this.state.uploading &&
          (
            <div>
              <ProgressBar now={uploadPercent} label={uploadBarLabel} />
              <h3>Uploading {this.state.currentUploadNum} of {this.state.uploadLength}</h3>
            </div>
          )
        }

        {this.state.uploadErrors.length > 0 && (
          this.state.uploadErrors.map(e => {
            return (
              <div>
                <p style={{ color: "red" }}>Error uploading {e.filename} : {e.error}</p>
              </div>
            )
          })
        )}

        <form enctype="multipart/form-data" style={{ textAlign: 'left' }} >
          <input type="file" name="files" onChange={this.onChange} multiple /><br /><br />
          <input type="button" value="upload" onClick={this.handleUpload.bind(this)} />
        </form>


        <InfiniteScroll
          dataLength={this.state.numPhotos}
          next={this.fetchMoreData}
          hasMore={this.state.hasMorePhotos}
          loader={<h4>Loading...</h4>}
        >
          {this.state.photosByMonth && this.state.photosByMonth.map(x => {
            if (!x.Photos) return;
            return (
              <div>
                <h2 align="left">{x.Month}</h2>
                <Gallery onClick={this.onClick} renderImage={this.imageRenderer} photos={
                  x.Photos.map(photo => {
                    return { id: photo.ID, src: "/photo_storage/thumbnails/" + photo.ThumbnailFilePath, height: 0, width: 0, photo: photo };
                  })
                } />
                <br />
              </div>
            )
          })}
        </InfiniteScroll>
      </div>
    )
  }
}


export default PhotosByMonthContainer;