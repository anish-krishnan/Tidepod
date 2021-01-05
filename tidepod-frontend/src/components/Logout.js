import React from 'react'
import { GoogleLogout } from 'react-google-login';

class Logout extends React.Component {
  onSuccess = (res) => {
    console.log('Logout made successfully');

    this.props.updateLoginStatus(false, '')
  };

  render() {
    return (
      <div>
        <GoogleLogout
          clientId={process.env.REACT_APP_GOOGLE_OAUTH_CLIENT_ID}
          buttonText="Logout"
          onLogoutSuccess={this.onSuccess}
        ></GoogleLogout>
      </div >
    )
  }
}

export default Logout;