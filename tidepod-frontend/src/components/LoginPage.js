import logo from '../sample_logo.png'
import tidepodLoginBackground from '../tidepodLoginBackground2.jpg'

import React from 'react'
import { FaLock } from 'react-icons/fa';

import Login from '../components/Login'

class LoginPage extends React.Component {
  render() {
    return (
      <div className="App" style={{ backgroundImage: `url(${tidepodLoginBackground})`, backgroundSize: 'cover', height: "100vh" }} >
        <br /><br /><br /><br /><br /><br />
        <h1 style={{ 'font-size': '60px' }}><img src={logo} className="App-logo" alt="logo" />Welcome to Tidepod</h1>

        <h3>{this.props.errorMessage}</h3>


        <div class="text-center text-success">
          <Login updateLoginStatus={this.props.updateLoginStatus} updateErrorStatus={this.props.updateErrorStatus} />
          <br /><br />
          <h1 class="text-center"><FaLock size="3em" /></h1>
        </div>

      </div >
    );
  }
}

export default LoginPage;