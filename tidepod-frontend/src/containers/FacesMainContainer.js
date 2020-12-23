import React from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";

class FacesMainContainer extends React.Component {
  state = {
    faces: []
  }

  componentDidMount() {
    fetch("http://localhost:3000/api/faces")
      .then(resp => resp.json())
      .then(faces => {
        console.log(faces)
        this.setState({
          faces: faces
        })
      })
  }



  render() {
    return (
      <div className="labels-main-container" >
        <h3>Faces Main Container</h3>

        {this.state.faces.map(face => {
          return (
            <div class="card" style={{ width: '18rem' }}>
              <div class="card-body">
                <h5 class="card-title"><Link class="nav-link" to={`/face/${face.ID}`} > {face.Name}</Link></h5>
                <h6 class="card-subtitle mb-2 text-muted">Card subtitle</h6>
                <p class="card-text">Some quick example text to build on the card title and make up the bulk of the card's content.</p>
              </div>
            </div>

          )
        })
        }
      </div >
    )
  }
}


export default FacesMainContainer;