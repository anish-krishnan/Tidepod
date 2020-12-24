import React from 'react'
import PhotoGallery from '../containers/PhotoGallery'

class LabelMainContainer extends React.Component {
  state = {
    label: "",
    photos: [],
    labelId: this.props.match.params.labelId
  }

  componentDidMount() {
    fetch("/api/label/" + this.state.labelId)
      .then(resp => resp.json())
      .then(label => {
        this.setState({
          label: label,
          photos: label.Photos
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

  render() {
    return (
      <div className="labels-main-container" >
        <h3>Label Main Container : {this.state.label.LabelName}</h3>
        <PhotoGallery photos={this.state.photos} handleDelete={this.handleDelete} />
      </div>
    )
  }
}


export default LabelMainContainer;