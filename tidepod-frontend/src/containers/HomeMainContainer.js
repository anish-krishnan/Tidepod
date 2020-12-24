import React from 'react'
import PhotoGallery from './PhotoGallery'

class HomeMainContainer extends React.Component {
  state = {
    photos: []
  }

  componentDidMount() {
    fetch("/api/photos")
      .then(resp => resp.json())
      .then(photos => {
        this.setState({
          photos: photos
        })
      })
  }

  handleDelete = (photo) => {
    fetch("/api/photos/delete/" + photo.ID,
      { method: "POST" })
      .then(resp => resp.json())
      .then(photos => {
        this.setState({
          photos: photos
        })
      })
  }

  handleUpload = () => {

  }

  render() {
    return (
      <div className="main-container" >
        <form enctype="multipart/form-data" action="/api/upload" method="post">
          <input type="file" name="files" multiple /><br /><br />
          <input type="submit" value="upload" />
        </form>
        <PhotoGallery photos={this.state.photos} handleDelete={this.handleDelete} />
      </div>
    )
  }
}


export default HomeMainContainer;