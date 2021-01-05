import React from 'react'
import { GoogleLogin } from 'react-google-login';
// refresh token
import { refreshTokenSetup } from '../utils/refreshToken';

class Login extends React.Component {
  onSuccess = (res) => {
    if (res.profileObj.googleId !== process.env.REACT_APP_VALID_GOOGLE_ID) {
      this.props.updateErrorStatus("Error: " + res.profileObj.email + " is not authorized")
      console.log("email address not in system")
      return
    }

    console.log('[Login Success] currentUser:', res);

    // initializing the setup
    refreshTokenSetup(res);

    if (res.accessToken) {
      this.props.updateLoginStatus(true, res.accessToken, res.tokenId)
    }
  };

  onFailure = (res) => {
    console.log('[Login failed] res:', res);
  };

  render() {
    return (
      <div>
        <GoogleLogin
          clientId={process.env.REACT_APP_GOOGLE_OAUTH_CLIENT_ID}
          buttonText="Sign in with Google"
          onSuccess={this.onSuccess}
          onFailure={this.onFailure}
          cookiePolicy={'single_host_origin'}
          style={{ marginTop: '100px' }}
          isSignedIn={true}
        />
      </div >
    )
  }
}

export default Login;