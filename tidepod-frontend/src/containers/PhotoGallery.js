import React from 'react'
import Photo from '../components/Photo'
import Gallery from 'react-photo-gallery';

class PhotoGallery extends React.Component {


  state = {
    isPhotoFullscreen: false,
    curPhoto: null
  };


  onClick = (event) => {
    this.setState({
      isPhotoFullscreen: !this.state.isPhotoFullscreen,
      curPhoto: { ...this.props.photos.find(photo => photo.ID === parseInt(event.target.id)) }
    });
    // debugger;
  }

  handleMinimize = () => {
    this.setState({
      isPhotoFullscreen: false
    })
  }

  handleDelete = (photo) => {
    this.props.handleDelete(photo)
    this.setState({
      isPhotoFullscreen: false
    })
    window.location.reload();
  }


  render() {
    if (this.state.isPhotoFullscreen) {
      return (<Photo photo={this.state.curPhoto} handleMinimize={this.handleMinimize} handleDelete={this.handleDelete} />)
    } else {
      return (
        <div>
          <h4>Gallery</h4>
          <Gallery onClick={this.onClick} photos={
            this.props.photos.map(photo => {
              return { id: photo.ID, src: "http://localhost:3000/photo_storage/thumbnails/" + photo.FilePath }
            })
          }
          />
        </div>
      )
    }
  }
}

export default PhotoGallery;