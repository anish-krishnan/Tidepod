import React, { useState } from 'react';
import {
  SafeAreaView,
  Text,
  Image,
  StyleSheet,
} from 'react-native';
import * as MediaLibrary from 'expo-media-library';
import * as Permissions from 'expo-permissions';

export default class UploadMainContainer extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      uploadLength: 0,
      currentUploadNum: 0,
      currentImageURI: "",
      uploading: false,

      errorFound: false,
      errorStatus: "",
    };
  }

  fetchAlbums = async () => {
    let { status } = await Permissions.askAsync(Permissions.MEDIA_LIBRARY);

    if (status !== 'granted') {
      console.log("Permissions not granted!")
    }
    let assetsPaged = await MediaLibrary.getAssetsAsync({ first: 10 });
    let assets = assetsPaged.assets

    // console.log("original assets\n", assetsPaged)

    let preformdata = new FormData()
    let formdata = new FormData()
    var asset, assetInfo;
    for (var i = 0; i < assets.length; i++) {
      asset = assets[i]
      assetInfo = await MediaLibrary.getAssetInfoAsync(asset)
      formdata.append('files', { uri: asset.uri, name: asset.filename })
      formdata.append('infoArray', JSON.stringify(assetInfo))
      preformdata.append('infoArray', JSON.stringify({ name: asset.filename, info: assetInfo }))
    }


    var photoIndices = await fetch("http://192.168.1.11:3000/api/preuploadmobile", {
      method: 'POST',
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      body: preformdata
    })
      .then(resp => resp.json())
      .catch(err => {
        console.log("error", err)
        this.setState({
          errorFound: true,
          errorStatus: err.toString()
        })
      })

    if (photoIndices == null) {
      this.setState({
        uploading: false
      })
      return;
    }

    var index;
    for (var i = 0; i < photoIndices.length; i++) {
      index = photoIndices[i]
      asset = assets[index]
      this.setState({
        uploading: true,
        currentImageURI: asset.uri,
        currentUploadNum: i + 1,
        uploadLength: photoIndices.length,
      })
      console.log("async", this.state.uploading, this.state.currentImageURI)
      assetInfo = await MediaLibrary.getAssetInfoAsync(asset)
      let filteredFormData = new FormData()
      filteredFormData.append('files', { uri: asset.uri, name: asset.filename })
      filteredFormData.append('infoArray', JSON.stringify(assetInfo))

      await fetch("http://192.168.1.11:3000/api/uploadmobile", {
        method: 'POST',
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        body: filteredFormData
      }).then(response => {
        console.log("response", response)
      }).catch(err => {
        console.log("error", err)
        this.setState({
          errorFound: true,
          errorStatus: err.toString()
        })
      })
    }
    this.setState({
      uploading: false
    })

  }

  componentDidMount() {
    this.fetchAlbums()
  }


  render() {
    if (this.state.errorFound) {
      return (
        <SafeAreaView style={styles.container}>
          <Text>Error: {this.state.errorStatus}</Text>
        </SafeAreaView>
      )
    } else if (this.state.uploading) {
      return (
        <SafeAreaView style={styles.container}>
          <Text>Uploading {this.state.currentUploadNum} of {this.state.uploadLength}</Text>
          <Image source={{ width: 200, height: 200, uri: this.state.currentImageURI }} />
        </SafeAreaView>
      );
    } else {
      return (
        <SafeAreaView style={styles.container}>
          <Text>Everything is up to date</Text>
        </SafeAreaView>
      );
    }
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
});