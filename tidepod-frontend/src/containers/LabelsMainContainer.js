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
import { FaTag } from 'react-icons/fa';

class LabelsMainContainer extends React.Component {
  state = {
    labels: []
  }

  componentDidMount() {
    fetch("/api/labels")
      .then(resp => resp.json())
      .then(labels => {
        if (labels)
          this.setState({
            labels: labels
          })
      })
  }



  render() {
    const labels = this.state.labels

    if (!labels.length) {
      return (
        <div className="labels-main-container" >
          <h3><FaTag /> No Labels</h3 >
          <br />
          <p>Upload more photos and come back later &#128512;</p>
        </div>
      )
    } else {
      return (
        <div className="labels-main-container" >
          <h3><FaTag /> Labels</h3>

          <Container>
            <Row>
              {labels.map(label => {
                const photo = label.Photos[0]
                return (
                  <Col xl={2} lg={3} md={4} sm={6} xs={12}>
                    <div class="card" style={{ maxWidth: '18rem' }}>
                      <img class="card-img-top" src={"/photo_storage/thumbnails/" + photo.FilePath} />
                      <div class="card-body">
                        <h5 class="card-title"><Link class="nav-link" to={`/label/${label.ID}`} > {label.LabelName} </Link></h5>
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


export default LabelsMainContainer;