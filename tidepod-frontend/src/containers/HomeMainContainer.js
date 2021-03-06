import React from 'react'
import Gallery from 'react-photo-gallery';

class HomeMainContainer extends React.Component {
  state = {
    photos: [],
    files: []
  }

  componentDidMount() {
    fetch("/api/photos", {
      headers: { "Token": this.props.idToken }
    })
      .then(resp => resp.json())
      .then(photos => {
        this.setState({
          photos: photos
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
        <Gallery onClick={this.onClick} photos={
          this.state.photos.map(photo => {
            return { id: photo.ID, src: "/photo_storage/thumbnails/" + photo.ThumbnailFilePath, height: 0, width: 0 }
          })
        } />
      </div>
    )
  }
}


export default HomeMainContainer;