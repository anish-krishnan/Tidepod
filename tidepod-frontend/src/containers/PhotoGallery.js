import React from 'react'
import Photo from '../components/Photo'
import Gallery from 'react-photo-gallery';
import { useHistory, withRouter } from "react-router-dom";

class PhotoGallery extends React.Component {


  state = {
    isPhotoFullscreen: false,
    curPhoto: null
  };


  // onClick = (event) => {
  //   this.setState({
  //     isPhotoFullscreen: !this.state.isPhotoFullscreen,
  //     curPhoto: { ...this.props.photos.find(photo => photo.ID === parseInt(event.target.id)) }
  //   });
  //   // debugger;
  // }

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

  onClick = (event) => {
    debugger;
    const { history } = this.props;
    if (history) history.push('/photo/' + event.target.id);
  }


  render() {
    if (this.state.isPhotoFullscreen) {
      return (<Photo photo={this.state.curPhoto} handleMinimize={this.handleMinimize} handleDelete={this.handleDelete} />)
    } else {
      return (
        <div>
          <Gallery onClick={this.onClick} photos={
            this.props.photos.map(photo => {
              return { id: photo.ID, src: "/photo_storage/thumbnails/" + photo.FilePath, height: 0, width: 0 }
            })
          }
          />
        </div>
      )
    }
  }
}

export default PhotoGallery;