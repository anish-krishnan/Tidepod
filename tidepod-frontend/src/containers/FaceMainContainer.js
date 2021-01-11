import React from 'react'
import PhotoGallery from '../containers/PhotoGallery'
import Gallery from 'react-photo-gallery';
import { useHistory, withRouter } from "react-router-dom";
import { FaUser } from 'react-icons/fa';

class FaceMainContainer extends React.Component {
  state = {
    face: "",
    photos: [],
    faceId: this.props.match.params.faceId
  }

  componentDidMount() {
    fetch("/api/face/" + this.state.faceId)
      .then(resp => resp.json())
      .then(face => {
        this.setState({
          face: face,
          photos: face.Boxes.map(box => { return box.Photo })
        })
      })
  }

  onClick = (event) => {
    const { history } = this.props;
    if (history) history.push('/photo/' + event.target.id);
  }

  render() {
    return (
      <div className="labels-main-container" >
        <h1><FaUser /></h1>
        <h3>Photos of <b>{this.state.face.Name}</b></h3>
        <Gallery onClick={this.onClick} photos={
          this.state.photos.map(photo => {
            return { id: photo.ID, src: "/photo_storage/thumbnails/" + photo.FilePath, height: 0, width: 0 }
          })
        }
        />
      </div>
    )
  }
}


export default FaceMainContainer;