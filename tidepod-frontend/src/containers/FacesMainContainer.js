import React from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";

import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import { FaUsers } from 'react-icons/fa';

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
          <h3><FaUsers /> No Faces</h3>
          <br />
          <Link class="nav-link" to={"/unassignedBoxes"}>Tag more photos and come back later &#128512;</Link>
        </div>
      )
    } else {
      return (
        <div className="labels-main-container" >
          <h3><FaUsers /> Faces</h3>
          <Link class="nav-link" to={"/unassignedBoxes"}>Tag more photos</Link>
          <Container>
            <Row>
              {faces.map(face => {
                const box = face.Boxes[0]
                return (
                  <Col xl={2} lg={3} md={4} sm={6} xs={12}>
                    <div class="card" style={{ maxWidth: "18rem" }}>
                      <img class="card-img-top" src={"/photo_storage/boxes/" + box.FilePath} />
                      <div class="card-body">
                        <h4 class="card-title"><Link class="nav-link" to={`/face/${face.ID}`}> {face.Name}</Link></h4>
                      </div>
                    </div>
                  </Col>
                )
              })
              }
            </Row>
          </Container>
        </div >
      )
    }
  }
}


export default FacesMainContainer;