import logo from '../logo.svg';
import '../App.css';

import Nav from '../containers/Nav'
import HomeMainContainer from './HomeMainContainer'
import PhotoMainContainer from '../containers/PhotoMainContainer'
import LabelMainContainer from '../containers/LabelMainContainer'
import LabelsMainContainer from '../containers/LabelsMainContainer'
import Photo from '../components/Photo'
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
          <Route path='/' component={HomeMainContainer} exact />
          <Route path='/photo/:photoId' component={PhotoMainContainer} />
          <Route path='/labels' component={LabelsMainContainer} />
          <Route path='/label/:labelId' component={LabelMainContainer} />
          <Route component={Error} />
        </Switch>
      </Router>
    </div>
  );
}

export default App;
