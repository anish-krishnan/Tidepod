import React from 'react'
import PhotoGallery from './PhotoGallery'

class HomeMainContainer extends React.Component {
  state = {
    photos: []
  }

  componentDidMount() {
    fetch("http://localhost:3000/api/photos")
      .then(resp => resp.json())
      .then(photos => {
        this.setState({
          photos: photos
        })
      })
  }

  handleDelete = (photo) => {
    fetch("http://localhost:3000/api/photos/delete/" + photo.ID,
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
        <h3>Home Main Container</h3>
        <form enctype="multipart/form-data" action="http://localhost:3000/api/upload" method="post">
          <input type="file" name="files" multiple /><br /><br />
          <input type="submit" value="upload" />
        </form>
        <PhotoGallery photos={this.state.photos} handleDelete={this.handleDelete} />
      </div>
    )
  }
}


export default HomeMainContainer;