import React from 'react'
import Photo from '../components/Photo'

class PhotoMainContainer extends React.Component {
  state = {
    photo: "",
    photoId: this.props.match.params.photoId,
    fetching: true
  }

  componentDidMount() {
    fetch("/api/photo/" + this.state.photoId)
      .then(resp => resp.json())
      .then(photo => {
        console.log(photo)
        this.setState({
          photo: photo,
          fetching: false
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
    if (this.state.fetching) {
      return <p>Fetching...</p>
    } else {
      return (<Photo photo={this.state.photo} handleMinimize={this.handleMinimize} handleDelete={this.handleDelete} />)
    }
  }
}

export default PhotoMainContainer;