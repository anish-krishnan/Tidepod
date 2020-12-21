import React from 'react'

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

          {this.props.photo.Labels.map(function (label, i) {
            return (<p>LABELS: {label.LabelName}</p>);
          })}

          <button type="button" onClick={this.handleMinimize} className="btn btn-default">
            <span className="glyphicon glyphicon-resize-small" ></span>Minimize
          </button>
          <button type="button" onClick={this.handleDelete} className="btn btn-default">
            <span className="glyphicon glyphicon-trash"></span>Delete
          </button>

        </div>
      </div>
    )
  }
}

export default Photo;