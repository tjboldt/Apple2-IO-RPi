var fs = require('fs');
var gpio = require('onoff').Gpio;

var out_write = new gpio(5, 'high', {activeLow: true});
var out_read = new gpio(11, 'high', {activeLow: true});
var out_command1 = new gpio(9, 'low', {activeLow: false});
var out_command2 = new gpio(10, 'low', {activeLow: false});
var out_bit3 = new gpio(22, 'low', {activeLow: false});
var out_bit2 = new gpio(27, 'low', {activeLow: false});
var out_bit1 = new gpio(17, 'low', {activeLow: false});
var out_bit0 = new gpio(4, 'low', {activeLow: false});
var in_write = new gpio(12, 'in', 'both');
var in_read = new gpio(16, 'in', 'both');
var in_command1 = new gpio(20, 'in');
var in_command2 = new gpio(21, 'in');
var in_bit3 = new gpio(26, 'in');
var in_bit2 = new gpio(19, 'in');
var in_bit1 = new gpio(13, 'in');
var in_bit0 = new gpio(6, 'in');

function onCleanUp() {
  in_write.unwatchAll();
  in_write.unexport();
}

function onWriteChanged(err, value) {
  if (err) {
  }
  else {
    //console.log(`value: ${value}`); 
  }
}

number = 0;
startTime = Date.now();

function onReadChanged(err, value) {
  if (err) {
  }
  else {
    if (value == 1) {
      //console.log("Nibble has been read");
      out_write.writeSync(0);
    }
    else {
      writeNibble(number++);
      if(number >= 4096) {
        seconds = Date.now() - startTime; 
	console.log(`Sent 2 KiB in ${seconds} milliseconds.`);
        process.exit() 
      }
    }
  }
}

in_write.watch(onWriteChanged);
in_read.watch(onReadChanged);

function writeNibble(nibble) {
  out_bit3.writeSync((nibble&8)>>3);
  out_bit2.writeSync((nibble&4)>>2);
  out_bit1.writeSync((nibble&2)>>1);
  out_bit0.writeSync(nibble&1);
  out_command1.writeSync(0);
  out_command2.writeSync(0);
  out_read.writeSync(0);
  out_write.writeSync(1);
}

writeNibble(0);

process.on('SIGINT', onCleanUp);


