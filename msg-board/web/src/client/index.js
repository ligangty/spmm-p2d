'use strict'

import React from 'react';
import {render} from 'react-dom';

import MessageBoard from './view/MessageBoard.js';
import RegisterForm from './view/Register.js';
import { BrowserRouter as Router, Route } from "react-router-dom";
// import {  Route, Switch } from "react-router";

const Routing = () => {
  return (
    <Router>
      <div>
        <Route exact path="/" component={RegisterForm} />
        <Route exact path="/register" component={RegisterForm} />
        <Route exact path="/messages" component={MessageBoard} />
      </div>
    </Router>
  );
};

render(<Routing />, document.getElementById('root'));
