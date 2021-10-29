'use strict'

import React from 'react';
import { withRouter } from 'react-router' 

class RegisterForm extends React.Component{
  constructor(props){
    super(props);

    this.state = {
      username: '',
      mails:''
    };

    this.handleNameChange = this.handleNameChange.bind(this);
    this.handleEmailChange = this.handleEmailChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }


  handleNameChange(event) {
    this.setState({username: event.target.value});
  }

  handleEmailChange(event){
    this.setState({mails: event.target.value});
  }

  handleSubmit(event){
    event.preventDefault();
    let usernameVal = this.state.username;
    let mails = this.state.mails
    let url = `http://localhost:8080/user/${usernameVal}`
    if (mails.length>0) {
      url = url+`?mails=${mails}`
    }
    fetch(url, {
      method: "PUT",
      credentials: 'same-origin',
      headers: {
        "Content-Type": "application/json",
      }
    })
    .then(response => {
      if(response.ok){
        response.json().then(data=>{
          this.props.history.push({
            pathname: '/messages',
            state: {
              user: data
            }
          })
        });
      }
    });    
  }

  render(){
    return (
      <div>
        <form onSubmit={this.handleSubmit}>
        	<p>
            <label>Register a username:</label>
            <input id="username" name="username" type="text" value={this.state.username} onChange={this.handleNameChange}/><br />
          </p>
          <p>
            <label>User's emails:</label>
            <input id="mails" name="mails" type="text" value={this.state.mails} onChange={this.handleEmailChange} />
          </p>
        	<p>
            <input type="submit" />
            <button>Cancel</button>
          </p>
        </form>
      </div>
    );
  }
}

export default withRouter(RegisterForm)