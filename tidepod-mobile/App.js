import { StatusBar } from 'expo-status-bar';
import React, { Component } from 'react';
import { StyleSheet, Button, Text, Navigator, SafeAreaView } from 'react-native';
import HomeMainContainer from './containers/HomeMainContainer';
import UploadMainContainer from './containers/UploadMainContainer';
import * as Google from 'expo-google-app-auth';

export default class App extends Component {
  constructor(props) {
    super(props)
    this.state = {
      signedIn: false,
      name: "",
    }
  }

  signIn = async () => {
    try {
      const result = await Google.logInAsync({
        iosClientId: "490899834821-6jviap0vgmp38l0huandpvb9hnf0ib9e.apps.googleusercontent.com",
        scopes: ['profile', 'email'],
      });

      if (result.type === 'success') {
        console.log("Login successful", result)
        this.setState({
          signedIn: true
        })
        return result.accessToken;
      } else {
        console.log("Login cancelled")
        return { cancelled: true };
      }
    } catch (e) {
      return { error: true };
    }
  }

  render() {
    if (this.state.signedIn) {
      return (
        <UploadMainContainer />
      );
    } else {
      return (
        <SafeAreaView>
          <Text>Not logged in</Text>
          <Button title="Sign in with Google" onPress={() => this.signIn()} />
        </SafeAreaView>
      )
    }
  }
}