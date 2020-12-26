// Example of Grid Image Gallery in React Native
// https://aboutreact.com/grid-image-gallery/

// import React in our code
import React, { useState, useEffect } from 'react';

// import all the components we are going to use
import {
  SafeAreaView,
  StyleSheet,
  Text,
  View,
  Image,
  TouchableOpacity,
  FlatList,
  Modal,
} from 'react-native';

//import FastImage
import FastImage from 'react-native-fast-image';

const PhotoGallery = (props) => {
  const [imageID, setImageID] = useState(-1);
  const [
    modalVisibleStatus, setModalVisibleStatus
  ] = useState(false);
  const [dataSource, setDataSource] = useState([]);

  useEffect(() => {
    let items = props.photos
    setDataSource(items);
  }, []);

  const showModalFunction = (visible, imageURL) => {
    //handler to handle the click on image of Grid
    //and close button on modal
    setImageID(imageURL);
    setModalVisibleStatus(visible);
  };

  return (
    <SafeAreaView style={styles.container}>
      {modalVisibleStatus ? (
        <Modal
          transparent={false}
          animationType={'fade'}
          visible={modalVisibleStatus}
          onRequestClose={() => {
            showModalFunction(!modalVisibleStatus, -1);
          }}>
          <View style={styles.modelStyle}>
            <Image
              style={styles.fullImageStyle}
              source={{ uri: "http://192.168.1.11:3000/photo_storage/saved/" + imageID + ".jpg" }}

            />
            <TouchableOpacity
              activeOpacity={0.5}
              style={styles.closeButtonStyle}
              onPress={() => {
                showModalFunction(!modalVisibleStatus, -1);
              }}>
              <Image
                source={{
                  uri:
                    'https://raw.githubusercontent.com/AboutReact/sampleresource/master/close.png',
                }}
                style={{ width: 35, height: 35 }}
              />
            </TouchableOpacity>
          </View>
        </Modal>
      ) : (
          <View style={styles.container}>
            <Text style={styles.titleStyle}>
              Tidepod Photos
          </Text>
            <FlatList
              data={dataSource}
              renderItem={({ item }) => (
                <View style={styles.imageContainerStyle}>
                  <TouchableOpacity
                    key={item.id}
                    style={{ flex: 1 }}
                    onPress={() => {
                      showModalFunction(true, item.id);
                    }}>
                    <Image
                      style={styles.imageStyle}
                      source={{
                        uri: item.src,
                      }}
                    />
                  </TouchableOpacity>
                </View>
              )}
              //Setting the number of column
              numColumns={3}
              keyExtractor={(item, index) => index.toString()}
            />
          </View>
        )}
    </SafeAreaView>
  );
};
export default PhotoGallery;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#ffffff',
  },
  titleStyle: {
    padding: 16,
    fontSize: 20,
    color: 'white',
    backgroundColor: 'dodgerblue',
  },
  imageContainerStyle: {
    flex: 1,
    flexDirection: 'column',
    margin: 1,
  },
  imageStyle: {
    height: 120,
    width: '100%',
  },
  fullImageStyle: {
    justifyContent: 'center',
    alignItems: 'center',
    height: '100%',
    width: '98%',
    resizeMode: 'contain',
  },
  modelStyle: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'rgba(0,0,0,0.9)',
  },
  closeButtonStyle: {
    width: 25,
    height: 25,
    top: 50,
    right: 20,
    position: 'absolute',
  },
});