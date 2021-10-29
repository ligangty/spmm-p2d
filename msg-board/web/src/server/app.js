const express = require('express');
const os = require('os');
const Util = require('./Util.js');

const app = express();

let server = app.listen(4000, function () {
   var host = server.address().address
   var port = server.address().port

   console.log("Example app listening at http://%s:%s", host, port)
});

app.use(express.static('build'));
