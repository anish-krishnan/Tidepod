import React from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import Button from 'react-bootstrap/Button';

class Photo extends React.Component {

  handleMinimize = () => {
    this.props.handleMinimize()
  }

  handleDelete = () => {
    this.props.handleDelete(this.props.photo)
  }


  render() {
    return (
      <div className="photo">
        <img style={this.mystyle} src={"http://localhost:3000/photo_storage/saved/" + this.props.photo.FilePath} height="auto" width="50%"></img>
        <div className="panel-footer">
          <p>#{this.props.photo.ID}</p>
          <p>Camera Model {this.props.photo.CameraModel}</p>
          <p>Location: {this.props.photo.LocationString}</p>
          <p>Timestamp {this.props.photo.Timestamp}</p>
          <p>FocalLength {this.props.photo.FocalLength}</p>
          <p>Aperture {this.props.photo.ApertureFStop}</p>

          <div>
            {this.props.photo.Labels.map(function (label, i) {
              return (<Link to={`/label/${label.ID}`}><Button variant="secondary">{label.LabelName}</Button></Link>);
            })}
          </div>

          <div>
            {this.props.photo.Boxes.map(function (box, i) {
              return (
                <img class="img-thumbnail" height='100px' width='100px' src={"http://localhost:3000/photo_storage/boxes/" + box.ID + ".jpg"} ></img>

              );
            })}
          </div>

          <div>
            <button type="button" class="btn btn-warning" onClick={this.handleMinimize} >
              <span className="glyphicon glyphicon-resize-small" ></span>Minimize
          </button>
            <button type="button" class="btn btn-warning" onClick={this.handleDelete} >
              <span className="glyphicon glyphicon-trash"></span>Delete
          </button>
          </div>

        </div>
      </div >
    )
  }
}

export default Photo;