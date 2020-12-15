
class App extends React.Component {
  render() {
    if (this.loggedIn) {
      return (<LoggedIn />);
    } else {
      return (<LoggedIn />);
    }
  }
}

class Home extends React.Component {
  render() {
    return (
      <div className="container">
        <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
          <h1>Jokeish</h1>
          <p>A load of Dad jokes XD</p>
          <p>Sign in to get access </p>
          <a onClick={this.authenticate} className="btn btn-primary btn-lg btn-login btn-block">Sign In</a>
        </div>
      </div>
    )
  }
}

class LoggedIn extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      jokes: [],
      photos: []
    };

    this.serverRequest = this.serverRequest.bind(this);
  }

  serverRequest() {
    $.get("http://localhost:3000/api/jokes", res => {
      this.setState({
        jokes: res
      });
    });

    $.get("http://localhost:3000/api/photos", res => {
      this.setState({
        photos: res
      });
    });
  }

  componentDidMount() {
    this.serverRequest();
  }

  render() {
    return (
      <div className="container">
        <div className="col-lg-12">
          <br />
          <span className="pull-right"><a onClick={this.logout}>Log out</a></span>
          <h2>Tidepod</h2>
          <p>Let's feed you with some funny Jokes!!!</p>
          <div className="row">
            {this.state.jokes.map(function (joke, i) {
              return (<Joke key={i} joke={joke} />);
            })}
          </div>
          <div className="row">
            {this.state.photos.map(function (photo, i) {
              return (<Photo key={i} photo={photo} />);
            })}
          </div>
        </div>
        <CreateJokeForm />
        <UploadPhotoForm />
      </div>
    )
  }
}

class CreateJokeForm extends React.Component {
  constructor(props) {
    super(props);
    this.state = { value: '' };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({ value: event.target.value });
  }

  handleSubmit(event) {
    $.post("http://localhost:3000/api/jokes/create/" + this.state.value)
    window.location.reload();
    event.preventDefault();
  }

  render() {
    return (
      <form onSubmit={this.handleSubmit}>
        <label>
          Joke:
          <input type="text" name="joke" value={this.state.value} onChange={this.handleChange} />
        </label>
        <input type="submit" value="Submit" />
      </form>
    );
  }
}

class UploadPhotoForm extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <form enctype="multipart/form-data" action="http://localhost:3000/api/upload" method="post">
        Files: <input type="file" name="files" multiple /><br /><br />
        <input type="submit" value="Submit" />
      </form>
    )
  }
}

class Joke extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      liked: "",
      jokes: [],
    };
    this.like = this.like.bind(this);
    this.delete = this.delete.bind(this);
    this.serverRequest = this.serverRequest.bind(this);
  }

  like() {
    let joke = this.props.joke;
    this.serverRequest(joke)
  }

  delete() {
    console.log("ASLDJLKSAJD")
    let joke = this.props.joke;
    $.post(
      "http://localhost:3000/api/jokes/delete/" + joke.id,
      {},
      res => {
        console.log("res... ", res);
        this.props.jokes = res;
      }
    );
    window.location.reload();
  }

  serverRequest(joke) {
    $.post(
      "http://localhost:3000/api/jokes/like/" + joke.id,
      { like: 1 },
      res => {
        console.log("res... ", res);
        this.setState({ liked: "Liked!", jokes: res });
        this.props.jokes = res;
      }
    );
  }

  render() {
    return (
      <div className="col-xs-4">
        <div className="panel panel-default">
          <div className="panel-heading">#{this.props.joke.id} <span className="pull-right">{this.state.liked}</span></div>
          <div className="panel-body">
            {this.props.joke.joke}
          </div>
          <div className="panel-footer">
            {this.props.joke.likes} Likes &nbsp;
            <a onClick={this.like} className="btn btn-default">
              <span className="glyphicon glyphicon-thumbs-up"></span>
            </a>

            <a onClick={this.delete} className="btn btn-default">
              <span className="glyphicon glyphicon-trash"></span>
            </a>
          </div>
        </div>
      </div>
    )
  }
}

class Photo extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      liked: "",
      jokes: [],
    };
    this.delete = this.delete.bind(this);
  }

  delete() {
    console.log("ASLDJLKSAJD")
    let photo = this.props.photo;
    $.post(
      "http://localhost:3000/api/photos/delete/" + photo.id,
      {},
      res => {
        console.log("res... ", res);
        this.props.jokes = res;
      }
    );
    window.location.reload();
  }

  mystyle = {
    height: "200px",
    width: "200px",
    overflow: "hidden"
  };

  render() {
    return (
      <div className="col-xs-4">
        <div className="panel panel-default">
          <div className="panel-heading">#{this.props.photo.id} </div>
          <div className="panel-body">
            <img style={this.mystyle} src={"../../saved/" + this.props.photo.FilePath} ></img>
          </div>
          <div className="panel-footer">
            <a onClick={this.delete} className="btn btn-default">
              <span className="glyphicon glyphicon-trash"></span>
            </a>
          </div>
        </div>
      </div>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('app'));
