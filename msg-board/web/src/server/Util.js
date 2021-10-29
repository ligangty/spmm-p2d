'use strict'

module.exports={
  makeId: function(length){
    let root = Math.random().toString(36).replace(/[^a-z]+/g, '');
    // console.log(root);
    return length?root.substr(0, length):root.substr(0,6);
  }
}
