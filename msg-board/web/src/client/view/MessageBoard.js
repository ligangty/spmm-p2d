'use strict'
import {styles} from "../style/style.js";
import React from 'react';

export default class MessageBoard extends React.Component {
  constructor(props){
    super(props);

    this.state = {
      user: {},
      content: '',
      messages: []
    };

    this.handleSend = this.handleSend.bind(this);
    this.handleFetchMessages = this.handleFetchMessages.bind(this);
    this.handleContentChange = this.handleContentChange.bind(this);
  }

  componentDidMount(){
    this.setState({user: this.props.location.state.user})
    this.handleFetchMessages()
  }

  handleContentChange(event){
    this.setState({content: event.target.value});
  }

  handleFetchMessages() {
    fetch(`http://localhost:8080/message/all`, {
      method: "GET",
      credentials: 'same-origin',
      headers: {
        "Content-Type": "application/json",
      }
    })
    .then(response => {
      if(response.ok){
        response.json().then(data=>{
          this.setState({
            messages: data
          });
        });
      }
    })
  }

  handleSend(event){
    event.preventDefault();
    let newMessage = {
      "uid": this.state.user.id,
      "content": this.state.content
    }
    console.log(newMessage)
    fetch(`http://localhost:8080/message`, {
      method: "PUT",
      credentials: 'same-origin',
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(newMessage)
    })
    .then(response => {
      if(response.ok){
        this.handleFetchMessages()
        this.setState({content:""})
      }
    });    
  }

  render() {
    return (
			// use key to control re-render all component
      <div>
        <div>
          <h1>Hello, {this.state.user.name}</h1>
        </div>
        <div style={styles.InputContainer}>
          <form>
            <label htmlFor="message">message:</label>{" "}
            <textarea id="message" name="message" value={this.state.content} onChange={this.handleContentChange} rows="6" />{" "} 
            <button onClick={this.handleSend}>Send</button>
          </form>
        </div>
				<hr />
				<div>
          <table style={styles.MessageTable}>
            <thead>
              <tr style={{borderBottom: "1pt solid black"}}>
                <td>user</td>
                <td>content</td>
                <td>time</td>
              </tr>
            </thead>
            <tbody>
              { this.state.messages.map(function(msg, i){
                return (
                  <tr key={msg.id} style={{borderBottom: "1pt solid black"}}>
                    <td>{msg.uid}</td>
                    <td>{msg.content}</td>
                    <td>{msg.time}</td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      </div>
    );
  }
}