import React, { useCallback } from 'react'
import Gallery from 'react-photo-gallery';

import { FaPlayCircle } from 'react-icons/fa';


class PhotosByMonthContainer extends React.Component {
  state = {
    photosByMonth: [],
    files: []
  }

  componentDidMount() {
    fetch("/api/photosByMonth", {
      headers: { "Token": this.props.idToken }
    })
      .then(resp => resp.json())
      .then(photosByMonth => {
        console.log(photosByMonth)

        this.setState({
          photosByMonth: photosByMonth
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
      var file = this.state.files[i]
      const formData = new FormData();
      formData.append('files', file)
      formData.append('infoArray', file.lastModified)
      formData.append('token', this.props.idToken)
      await fetch("/api/upload", {
        method: 'POST',
        body: formData,
      })
        .then(resp => resp.json())
        .then(response => {
          console.log("response", response)
        })
        .catch(err => {
          console.log("error", err)
        })
    }
  }

  onClick = (event) => {
    const { history } = this.props;
    console.log("history", history)
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
    //   key={key}
    //   margin={"2px"}
    //   index={index}
    //   photo={photo}
    //   left={left}
    //   top={top}
    // />
  }

  // position: "absolute",
  // top: 0,
  // bottom: 0,
  // left: 0,
  // right: 0,
  // height: "100%",
  // width: "100%",
  // opacity: 0,
  // transition: ".3s ease",
  // "background-color": "red",

  render() {
    return (
      <div className="main-container" >
        <form enctype="multipart/form-data" >
          <input type="file" name="files" onChange={this.onChange} multiple /><br /><br />
          <input type="button" value="upload" onClick={this.handleUpload.bind(this)} />
        </form>

        {this.state.photosByMonth && this.state.photosByMonth.map(x => {
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
      </div>
    )
  }
}


export default PhotosByMonthContainer;