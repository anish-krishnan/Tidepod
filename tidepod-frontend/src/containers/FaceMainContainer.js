import React from 'react'
import PhotoGallery from '../containers/PhotoGallery'

class LabelMainContainer extends React.Component {
  state = {
    face: "",
    photos: [],
    faceId: this.props.match.params.faceId
  }

  componentDidMount() {
    fetch("http://localhost:3000/api/face/" + this.state.faceId)
      .then(resp => resp.json())
      .then(face => {
        this.setState({
          face: face,
          photos: face.Boxes.map(box => { return box.Photo })
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

  render() {
    return (
      <div className="labels-main-container" >
        <h3>Face Main Container : {this.state.face.Name}</h3>
        <PhotoGallery photos={this.state.photos} handleDelete={this.handleDelete} />
      </div>
    )
  }
}


export default LabelMainContainer;