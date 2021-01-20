
import React from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import Gallery from 'react-photo-gallery';
import Container from 'react-bootstrap/Container'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'
import { FaSearch } from 'react-icons/fa';

class SearchMainContainer extends React.Component {

  onClick = (event) => {
    const { history } = this.props;
    if (history) history.push('/photo/' + event.target.id);
  }

  render() {
    const query = this.props.location.state.query
    const labels = this.props.location.state.result.Labels
    const faces = this.props.location.state.result.Faces
    const photos = this.props.location.state.result.Photos

    return (
      <div className="search-main-container" >
        <h3><FaSearch /> Results for "{query}"</h3>
        <br />

        {labels && (
          <div>
            <h4>Related Labels</h4>
            {labels.map((label, i) => {
              return (<Link class="badge badge-secondary" style={{ "margin-right": "5px" }} to={`/label/${label.ID}`} key={i}><h5>{label.LabelName}</h5></Link>);
            })}
            <br /><br />
          </div>
        )}


        {faces && (
          <div>
            <h4>Related Faces</h4>
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
            <br /><br />
          </div>
        )}

        {photos && (<Gallery onClick={this.onClick} photos={
          photos.map(photo => {
            return { id: photo.ID, src: "/photo_storage/thumbnails/" + photo.FilePath, height: 0, width: 0 }
          })
        }
        />)}

        {!photos && (<h3>No Results :(</h3>)}

      </div >
    )
  }
}


export default SearchMainContainer;