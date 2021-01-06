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
    fetch("/api/faces")
      .then(resp => resp.json())
      .then(faces => {
        console.log(faces)
        this.setState({
          faces: faces
        })
      })
  }



  render() {
    const faces = this.state.faces
    if (!faces.length) {
      return (
        <div className="labels-main-container" >
          <h3>No Faces</h3>
          <br />
          <p>Tag more photos and come back later &#128512;</p>
        </div>
      )
    } else {
      return (
        <div className="labels-main-container" >
          <h3>Faces</h3>

          <div class="card-deck">
            {faces.map(face => {
              const box = face.Boxes[0]
              return (
                <div class="card" style={{ maxWidth: "200px" }}>
                  <img class="card-img-top" src={"/photo_storage/boxes/" + box.FilePath} />
                  <div class="card-body">
                    <h4 class="card-title"><Link class="nav-link" to={`/face/${face.ID}`} > {face.Name}</Link></h4>
                  </div>
                </div>
              )
            })
            }
          </div>
        </div >
      )
    }
  }
}


export default FacesMainContainer;