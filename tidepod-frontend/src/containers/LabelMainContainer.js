import React from 'react'
import PhotoGallery from '../containers/PhotoGallery'
import Gallery from 'react-photo-gallery';
import { FaTag } from 'react-icons/fa';

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

  onClick = (event) => {
    const { history } = this.props;
    if (history) history.push('/photo/' + event.target.id);
  }

  render() {
    return (
      <div className="labels-main-container" >
        <h1><FaTag /></h1>
        <h3>Photos of a {this.state.label.LabelName}</h3>
        <Gallery onClick={this.onClick} photos={
          this.state.photos.map(photo => {
            return { id: photo.ID, src: "/photo_storage/thumbnails/" + photo.FilePath, height: 0, width: 0 }
          })
        } />
      </div>
    )
  }
}


export default LabelMainContainer;