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
import { FaArrowLeft } from 'react-icons/fa';
mapboxgl.accessToken = process.env.REACT_APP_MAPBOX_API_KEY;

class Photo extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      lng: this.props.photo.Longitude,
      lat: this.props.photo.Latitude,
      zoom: 16
    };
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

  handleMinimize = () => {
    this.props.handleMinimize()
  }

  handleDelete = () => {
    this.props.handleDelete(this.props.photo)
  }

  render() {
    return (

      <div className="photo" style={{ 'backgroundColor': 'black' }}>
        <div class="container-fluid" >
          <div class="row" >
            <div class="col-9" >
              <a href="/"><h1 style={{ 'color': 'white', 'position': 'absolute', 'top': 0, 'left': 0 }}><FaArrowLeft /></h1></a>
              <img style={{ 'height': 'auto', 'maxWidth': '100%', 'maxHeight': '95vh' }} src={"/photo_storage/saved/" + this.props.photo.FilePath} ></img>
            </div>

            <div class="col-3 float-right" style={{ 'backgroundColor': 'white' }} align="left">
              <h3>Info</h3>
              <table class="table table-bordered">
                <tbody>
                  <tr>
                    <td><b>Camera Model</b></td>
                    <td>{this.props.photo.CameraModel}</td>
                  </tr>
                  <tr>
                    <td><b>Location</b></td>
                    <td>{this.props.photo.LocationString}</td>
                  </tr>
                  <tr>
                    <td><b>Timestamp</b></td>
                    <td>{this.props.photo.Timestamp}</td>
                  </tr>
                  <tr>
                    <td><b>FocalLength</b></td>
                    <td>{this.props.photo.FocalLength}</td>
                  </tr>
                  <tr>
                    <td><b>Aperture</b></td>
                    <td>{this.props.photo.ApertureFStop}</td>
                  </tr>
                </tbody>
              </table>

              {this.props.photo.LocationString.length > 0 &&
                <div ref={el => this.mapContainer = el} style={{
                  'width': '100%', 'height': '40%',
                }} align="center" />
              }

              <div>
                {this.props.photo.Labels.map(function (label, i) {
                  return (<Link to={`/label/${label.ID}`} key={i}><Button variant="secondary">{label.LabelName}</Button></Link>);
                })}
              </div>

              <div>
                {this.props.photo.Boxes.map(function (box, i) {
                  console.log(box)
                  return (<Box box={box} key={i} />);
                })}
              </div>

              <div>
                <button type="button" class="btn btn-warning" onClick={this.handleDelete} >
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

export default Photo;