import React from 'react'
import { GoogleLogin } from 'react-google-login';
// refresh token
import { refreshTokenSetup } from '../utils/refreshToken';

const clientId = '490899834821-99nr82qe6u0719iub1ojtuj2ijc7vhp4.apps.googleusercontent.com'

class Login extends React.Component {
  onSuccess = (res) => {
    if (res.profileObj.email !== "anishtech1@gmail.com") {
      console.log("email address not in system")
      return
    }

    console.log('[Login Success] currentUser:', res.profileObj);

    // initializing the setup
    refreshTokenSetup(res);

    if (res.accessToken) {
      this.props.updateLoginStatus(true, res.accessToken)
    }
  };

  onFailure = (res) => {
    console.log('[Login failed] res:', res);
  };

  render() {
    return (
      <div>
        <GoogleLogin
          clientId={clientId}
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