import React from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";
import Button from 'react-bootstrap/Button';
import Box from '../components/Box'
import mapboxgl from 'mapbox-gl';
import { FaArrowLeft, FaDownload } from 'react-icons/fa';
import { withRouter } from 'react-router-dom'
import ReactPlayer from 'react-player/file';
mapboxgl.accessToken = process.env.REACT_APP_MAPBOX_API_KEY;

class Photo extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      lng: this.props.photo.Longitude,
      lat: this.props.photo.Latitude,
      zoom: 16
    };
    this.goBack = this.goBack.bind(this);
  }

  componentDidMount() {
    if (this.props.photo.LocationString.length == 0) {
      return;
    }
    const map = new mapboxgl.Map({
      container: this.mapContainer,
      style: 'mapbox://styles/mapbox/streets-v11',
      center: [this.state.lng, this.state.lat],
      zoom: this.state.zoom
    });

    const marker = new mapboxgl.Marker()
      .setLngLat([this.state.lng, this.state.lat])
      .addTo(map);
  }

  formatDate = (string) => {
    var options = { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long', hour: '2-digit', minute: '2-digit', second: '2-digit', timeZoneName: 'short' };
    return new Date(string).toLocaleDateString([], options);
  }

  goBack = (event) => {
    this.props.history.goBack()
  }

  handleMinimize = () => {
    this.props.handleMinimize()
  }

  handleDelete = () => {
    this.props.handleDelete(this.props.photo)
  }

  render() {
    const photo = this.props.photo
    return (

      <div className="photo" >
        <div class="container-fluid" >
          <div class="row" >
            <div class="col-9" style={{ 'backgroundColor': 'black', height: "100vh" }}>
              {
                photo.MediaType == "photo" && (<img style={{ 'height': 'auto', 'maxWidth': '100%', 'maxHeight': '95vh' }} src={"/photo_storage/saved/" + photo.FilePath} />)
              }
              {
                photo.MediaType == "video" && (<ReactPlayer url={"/photo_storage/saved/" + photo.FilePath} width="100%" height="90%" controls={true} />)
              }
              <a onClick={this.goBack} ><h1 style={{ 'color': 'white', 'position': 'absolute', 'top': 0, 'left': 0 }}><FaArrowLeft /></h1></a>
            </div>

            <div class="col-3 float-right" style={{ 'backgroundColor': 'white' }} align="left">
              <h3>Info</h3>
              <table class="table table-bordered">
                <tbody>
                  <tr>
                    <td><b>Camera Model</b></td>
                    <td>{photo.CameraModel}</td>
                  </tr>
                  <tr>
                    <td><b>Location</b></td>
                    <td>{photo.LocationString}</td>
                  </tr>
                  <tr>
                    <td><b>Timestamp</b></td>
                    <td>{this.formatDate(photo.Timestamp)}</td>
                  </tr>
                  <tr>
                    <td><b>FocalLength</b></td>
                    <td>{photo.FocalLength}</td>
                  </tr>
                  <tr>
                    <td><b>Aperture</b></td>
                    <td>{photo.ApertureFStop}</td>
                  </tr>
                  <tr>
                    <td><b>Filename</b></td>
                    <td>{photo.OriginalFilename}</td>
                  </tr>
                </tbody>
              </table>

              <div>
                {photo.Labels.map(function (label, i) {
                  return (<Link class="badge badge-secondary" style={{ "margin-right": "5px" }} to={`/label/${label.ID}`} key={i}><h5>{label.LabelName}</h5></Link>);
                })}
              </div>
              <br />

              <div class="card-deck">
                {photo.Boxes.map(function (box, i) {
                  console.log(box)
                  return (<Box box={box} key={i} />);
                })}
              </div>
              <br />

              {photo.LocationString.length > 0 &&
                <div ref={el => this.mapContainer = el} style={{
                  'width': '100%', 'height': '40%',
                }} align="center" />
              }
              <br />

              <div>
                <a href={"http://localhost:3000/photo_storage/saved/" + photo.FilePath} style={{ 'color': 'black' }} download><FaDownload /></a>
                <button type="button" class="btn btn-danger" onClick={this.handleDelete} >
                  <span className="glyphicon glyphicon-trash"></span>Delete
                </button>
              </div>
            </div>
          </div>
        </div>
      </div >
    )
  }
}

export default withRouter(Photo);