import React, { Component } from 'react';
import { StyleSheet, Text, TouchableHighlight, Button, View, Alert, Image, SafeAreaView } from 'react-native';
import PhotoGallery from '../containers/PhotoGallery'

export default class HomeMainContainer extends Component {

  constructor(props) {
    super(props);

    this.state = {
      photos: [],
      loading: true,
    };
  }

  componentDidMount() {
    fetch("http://192.168.1.11:3000/api/photos")
      .then(resp => resp.json())
      .then(photos => {
        this.setState({
          photos: photos,
          loading: false,
        })
      })
  }

  render() {
    console.log("rendering... loading=", this.state.loading)
    if (this.state.loading) {
      return (
        <SafeAreaView>
          <Text>Loading...</Text>
        </SafeAreaView>
      )
    } else {
      return (
        <PhotoGallery photos={
          this.state.photos.map(photo => {
            return { id: photo.ID, src: "http://192.168.1.11:3000/photo_storage/thumbnails/" + photo.FilePath, height: 0, width: 0 }
          })}
        />
      )
    }
  }
}