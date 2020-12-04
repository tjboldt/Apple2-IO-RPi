var fs = require('fs');
var gpio = require('onoff').Gpio;

var output = new gpio(4, 'high', {activeLow: true});
var input = new gpio(6, 'in', 'both');

function onCleanUp() {
  output.unwatchAll();
  output.unexport();
  input.unexport();
}

function onStatusChanged(err, value) {
  if (err) {
  }
  else {
  }
}

input.watch(onStatusChanged);

process.on('SIGINT', onCleanUp);
