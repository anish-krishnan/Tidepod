import React from 'react'
import Photo from '../components/Photo'

class PhotoMainContainer extends React.Component {

  state = {
    photo: "",
    photoId: this.props.match.params.photoId,
    fetching: true
  }

  componentDidMount() {
    fetch("/api/photo/" + this.state.photoId, {
      headers: { "Token": this.props.idToken }
    })
      .then(resp => resp.json())
      .then(photo => {
        console.log("GOT PHOTO", photo)
        this.setState({
          photo: photo,
          fetching: false
        })
      })
  }

  handleDelete = (photo) => {
    fetch("/api/photos/delete/" + photo.ID,
      {
        method: "POST",
        headers: { "Token": this.props.idToken }
      })
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