import React from 'react'
import PhotoGallery from './PhotoGallery'
import Gallery from 'react-photo-gallery';

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

  onClick = (event) => {
    const { history } = this.props;
    if (history) history.push('/photo/' + event.target.id);
  }

  render() {
    return (
      <div className="main-container" >
        <form enctype="multipart/form-data" action="/api/upload" method="post">
          <input type="file" name="files" multiple /><br /><br />
          <input type="submit" value="upload" />
        </form>
        <Gallery onClick={this.onClick} photos={
          this.state.photos.map(photo => {
            return { id: photo.ID, src: "/photo_storage/thumbnails/" + photo.FilePath, height: 0, width: 0 }
          })
        } />
      </div>
    )
  }
}


export default HomeMainContainer;