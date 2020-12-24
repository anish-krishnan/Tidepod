import React from 'react'

class Box extends React.Component {

  state = {
    input: '',
    box: this.props.box
  }

  handleListInput = (event) => {
    this.setState({
      input: event.target.value
    })
  }

  handleListSubmit = (event) => {
    event.preventDefault()

    console.log("SUBMITTED NAME", this.state.input)

    fetch("/api/boxes/assignface/" + this.state.box.ID + "+" + this.state.input)
      .then(resp => resp.json())
      .then(box => {
        this.setState({
          box: box
        })
      })
    // const requestOptions = {
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify({ name: 'React POST Request Example' })
    // };
    // fetch('http://localhost:3000/api/boxes/assignface/6', requestOptions)
  }


  render() {
    if (this.state.box.Face.ID == 0) {
      return (
        <div class="card" style={{ width: '10rem' }}>
          <img class="img-thumbnail" height='100px' width='100px' src={"/photo_storage/boxes/" + this.state.box.FilePath} />
          <form onSubmit={this.handleListSubmit}>

            <input onChange={this.handleListInput} type="text" class="form-control" placeholder="Enter name" />

            <button type="submit" class="btn btn-primary">Submit</button>
          </form>
        </div>
      )

    } else {
      return (
        <div class="card" style={{ width: '6rem' }}>
          <img class="img-thumbnail" height='100px' width='100px' src={"/photo_storage/boxes/" + this.state.box.FilePath} />
          <h5 class="card-title">{this.state.box.Face.Name}</h5>
        </div>
      )
    }
  }
}

export default Box
