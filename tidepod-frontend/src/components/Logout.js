import React from 'react'
import { GoogleLogout } from 'react-google-login';

const clientId = '490899834821-99nr82qe6u0719iub1ojtuj2ijc7vhp4.apps.googleusercontent.com'

class Logout extends React.Component {
  onSuccess = (res) => {
    console.log('Logout made successfully');

    this.props.updateLoginStatus(false, '')
  };

  render() {
    return (
      <div>
        <GoogleLogout
          clientId={clientId}
          buttonText="Logout"
          onLogoutSuccess={this.onSuccess}
        ></GoogleLogout>
      </div >
    )
  }
}

export default Logout;