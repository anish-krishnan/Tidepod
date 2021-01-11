import logo from '../sample_logo.png'
import tidepodLoginBackground from '../tidepodLoginBackground2.jpg'
import '../App.css';

import React from 'react'
import Nav from '../containers/Nav'
import HomeMainContainer from './HomeMainContainer'
import PhotosByMonthContainer from '../containers/PhotosByMonthContainer'
import PhotoMainContainer from '../containers/PhotoMainContainer'
import Photo from '../components/Photo'
import LabelMainContainer from '../containers/LabelMainContainer'
import LabelsMainContainer from '../containers/LabelsMainContainer'
import FacesMainContainer from '../containers/FacesMainContainer'
import FaceMainContainer from '../containers/FaceMainContainer'
import UnassignedBoxesContainer from '../containers/UnassignedBoxesContainer'
import LoginPage from '../components/LoginPage'

import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";



class App extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      isLoggedIn: false,
      accessToken: '',
      idToken: '',
      errorMessage: ''
    };
  }

  updateLoginStatus = (status, accessToken, idToken) => {
    this.setState({
      isLoggedIn: status,
      accessToken: accessToken,
      idToken: idToken,
      errorMessage: ''
    })
  }

  updateErrorStatus = (message) => {
    this.setState({
      errorMessage: message
    })
  }

  render() {

    if (this.state.isLoggedIn) {
      return (
        <div className="App" >
          <Router>
            <Nav updateLoginStatus={this.updateLoginStatus} />

            <Switch>
              <Route exact path='/' component={(props) => <PhotosByMonthContainer idToken={this.state.idToken} {...props} />} />

              <Route exact path='/allPhotos' component={(props) => <HomeMainContainer idToken={this.state.idToken} {...props} />} />
              <Route path='/photo/:photoId' component={(props) => <PhotoMainContainer idToken={this.state.idToken} {...props} />} />

              <Route path='/labels' component={LabelsMainContainer} />
              <Route path='/label/:labelId' component={LabelMainContainer} />

              <Route path='/faces' component={FacesMainContainer} />
              <Route path='/face/:faceId' component={FaceMainContainer} />

              <Route path='/unassignedBoxes' component={UnassignedBoxesContainer} />

              <Route component={Error} />
            </Switch>
          </Router>
        </div>
      );
    } else {
      return (
        <LoginPage errorMessage={this.state.errorMessage} updateLoginStatus={this.updateLoginStatus} updateErrorStatus={this.updateErrorStatus} />
      );
    }
  }
}

export default App;
