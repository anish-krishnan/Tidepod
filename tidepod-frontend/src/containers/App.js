import logo from '../logo.svg';
import '../App.css';

import Nav from '../containers/Nav'
import HomeMainContainer from './HomeMainContainer'
import PhotosByMonthContainer from '../containers/PhotosByMonthContainer'
import PhotoMainContainer from '../containers/PhotoMainContainer'
import Photo from '../components/Photo'
import LabelMainContainer from '../containers/LabelMainContainer'
import LabelsMainContainer from '../containers/LabelsMainContainer'
import FacesMainContainer from '../containers/FacesMainContainer'
import FaceMainContainer from '../containers/FaceMainContainer'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";

// {/* <Nav /> */}
function App() {
  return (
    <div className="App">
      <Router>
        <Nav />

        <Switch>
          <Route path='/' component={PhotosByMonthContainer} exact />
          <Route path='/byMonth' component={PhotosByMonthContainer} exact />
          <Route path='/photo/:photoId' component={PhotoMainContainer} />
          <Route path='/labels' component={LabelsMainContainer} />
          <Route path='/label/:labelId' component={LabelMainContainer} />
          <Route path='/faces' component={FacesMainContainer} />
          <Route path='/face/:faceId' component={FaceMainContainer} />
          <Route component={Error} />
        </Switch>
      </Router>
    </div>
  );
}

export default App;
