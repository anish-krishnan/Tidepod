import React from 'react'
import logo from '../sample_logo.png'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";

class Nav extends React.Component {

  handleRunFaceDetect = () => {
    fetch("/api/classifyFaces")
  }

  render() {
    return (
      <nav class="navbar navbar-expand-lg navbar-light bg-light">
        <a class="navbar-brand mb-0 h1" ><img src={logo} width="30" height="30" class="d-inline-block align-top" alt="" />Tidepod</a>
        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>

        <div class="collapse navbar-collapse" id="navbarSupportedContent">
          <ul class="navbar-nav mr-auto">
            <li class="nav-item active">
              <Link class="nav-link" to="/">Home<span class="sr-only">(current)</span></Link>
            </li>
            <li class="nav-item">
              <Link class="nav-link" to="/labels">Labels</Link>
            </li>
            <li class="nav-item">
              <Link class="nav-link" to="/faces">Faces</Link>
            </li>

            <button onClick={this.handleRunFaceDetect} class="btn btn-sm btn-outline-secondary" type="button">Run Face Detect</button>

          </ul>
        </div>
        <form class="form-inline my-2 my-lg-0">
          <input class="form-control mr-sm-2" type="search" placeholder="Search" aria-label="Search" />
          <button class="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>
        </form>
      </nav>
    )
  }
}
export default Nav;