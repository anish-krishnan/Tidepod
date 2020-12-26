import { StatusBar } from 'expo-status-bar';
import React from 'react';
import { StyleSheet, Text, TouchableHighlight, Button, View, Alert, Image, SafeAreaView } from 'react-native';
import HomeMainContainer from './containers/HomeMainContainer'

export default function App() {

  return (
    // <SafeAreaView style={styles.container}>
    //   <TouchableHighlight onPress={() => console.log("Image Pressed!")} onLongPress={() => console.log("Image Loooong press!")}>
    //     <Image source={{
    //       width: 100,
    //       height: 150,
    //       uri: "https://picsum.photos/200/300"
    //     }} />
    //   </TouchableHighlight>
    //   <Button
    //     title="Click Me"
    //     onPress={() =>
    //       Alert.alert("My Title", "Message", [
    //         { text: "Yes", onPress: () => console.log("Yes") },
    //         { text: "No", onPress: () => console.log("No") }
    //       ])
    //     }
    //   />
    //   <Button
    //     title="Click Me too"
    //     onPress={() =>
    //       Alert.prompt("My Title", "Message", text => console.log(text))
    //     }
    //   />
    // </SafeAreaView >
    <HomeMainContainer />
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
});
