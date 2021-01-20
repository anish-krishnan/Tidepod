import React, { useState } from 'react'
import logo from '../sample_logo.png'
import Logout from '../components/Logout'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link,
  useHistory
} from "react-router-dom";

function Nav(props) {
  const history = useHistory();

  let state = {
    query: '',
  }

  const [query, setQuery] = useState("");
  const [result, setResult] = useState(null);


  let handleRunFaceDetect = () => {
    fetch("/api/classifyFaces")
  }

  let handleChangeQuery = (event) => {
    setQuery(event.target.value)
  }

  let handleSearch = (event) => {
    event.preventDefault();

    const formattedQuery = query.split(" ").join("&")
    fetch("/api/search/" + formattedQuery)
      .then(resp => resp.json())
      .then(res => {
        console.log("GOT RESULT", res)

        history.push({
          pathname: '/search',
          state: { query: query, result: res }
        });

      })


  }

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

          <button onClick={handleRunFaceDetect} class="btn btn-sm btn-outline-secondary" type="button">Run Face Detect</button>
          <Logout updateLoginStatus={props.updateLoginStatus} />
        </ul>
      </div>
      <form onSubmit={handleSearch} class="form-inline my-2 my-lg-0">
        <input onChange={handleChangeQuery} class="form-control mr-sm-2" type="search" placeholder="Search" aria-label="Search" />
        <button class="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>
      </form>
    </nav>
  )

}
export default Nav;