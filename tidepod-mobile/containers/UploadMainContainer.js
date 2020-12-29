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

  fetchAlbums = async () => {
    let { status } = await Permissions.askAsync(Permissions.MEDIA_LIBRARY);

    if (status !== 'granted') {
      console.log("Permissions not granted!")
    }
    let assetsPaged = await MediaLibrary.getAssetsAsync({ first: 100 });
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



    var photoIndices = await fetch("http://73.71.1.40:3000/api/preuploadmobile", {
      method: 'POST',
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      body: preformdata
    })
      .then(resp => resp.json())
      .catch(err => {
        console.log("error", err)
      })

    var index;
    for (var i = 0; i < photoIndices.length; i++) {
      index = photoIndices[i]
      asset = assets[index]
      assetInfo = await MediaLibrary.getAssetInfoAsync(asset)
      let filteredFormData = new FormData()
      filteredFormData.append('files', { uri: asset.uri, name: asset.filename })
      filteredFormData.append('infoArray', JSON.stringify(assetInfo))
      await fetch("http://73.71.1.40:3000/api/uploadmobile", {
        method: 'POST',
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        body: filteredFormData
      }).then(response => {
        console.log("response", response)
      }).catch(err => {
        console.log("error", err)
      })
    }

  }

  componentDidMount() {
    this.fetchAlbums()
  }


  render() {
    return (
      <SafeAreaView style={styles.container}>
        <Text>This is the uploader</Text>
        <Image source={{ width: 100, height: 100, uri: "assets-library://asset/asset.JPG?id=202EE2C2-1397-49D0-A8E1-C75B7BEDD497&ext=JPG" }} />
      </SafeAreaView>
    );
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