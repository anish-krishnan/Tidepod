import React from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import Box from '../components/Box'

import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import { FaUserTag } from 'react-icons/fa';

class UnassignedBoxesContainer extends React.Component {
  state = {
    boxes: []
  }

  componentDidMount() {
    fetch("/api/unassignedBoxes")
      .then(resp => resp.json())
      .then(boxes => {
        this.setState({
          boxes: boxes
        })
      })
  }
  render() {
    const boxes = this.state.boxes

    if (!boxes || boxes.length == 0) {
      return (
        <div>
          <h3>Unassigned Boxes</h3>
          <br />
          <p>Everyone is tagged &#128512;</p>
        </div>
      )
    } else {
      return (
        <div>
          <h3><FaUserTag /> Unassigned Boxes</h3>
          <Container>
            <Row>
              {boxes.map(box => {
                return (
                  <Col xl={2} lg={3} md={4} sm={6} xs={12}>
                    <Box box={box} showExpand={true} />
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
export default UnassignedBoxesContainer;