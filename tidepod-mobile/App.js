import { StatusBar } from 'expo-status-bar';
import React, { Component } from 'react';
import { StyleSheet, Text, Navigator } from 'react-native';
import HomeMainContainer from './containers/HomeMainContainer';
import UploadMainContainer from './containers/UploadMainContainer';


export default class App extends Component {

  render() {
    console.log("inside render method")
    return (
      <UploadMainContainer />
    );
  }
}