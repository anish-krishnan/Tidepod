import React from 'react'
import Gallery from 'react-photo-gallery';


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
              <Gallery onClick={this.onClick} photos={
                x.Photos.map(photo => {
                  return { id: photo.ID, src: "/photo_storage/thumbnails/" + photo.ThumbnailFilePath, height: 0, width: 0 }
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