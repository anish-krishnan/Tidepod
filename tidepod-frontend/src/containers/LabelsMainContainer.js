import React from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";

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
          <h3>No Labels</h3 >
          <br />
          <p>Upload more photos and come back later &#128512;</p>
        </div>
      )
    } else {
      return (
        <div className="labels-main-container" >
          <h3>Labels</h3>
          {labels.map((label) => {
            return (
              <div class="card" style={{ width: '18rem' }}>
                <div class="card-body">
                  <h5 class="card-title"><Link class="nav-link" to={`/label/${label.ID}`} > {label.LabelName} : {label.Photos.length}</Link></h5>
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
}


export default LabelsMainContainer;