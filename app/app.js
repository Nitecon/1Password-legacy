var path = require('path');
var gui = require('nw.gui');
var spawn = require('child_process').spawn;
//var nwShell = gui.Shell;

// Main node entry point, we use this to launch the server end etc.
//var arch = process.arch;
var platform = process.platform;

var isWin = false;


var execPath = path.dirname(process.execPath);

// some base path is appended
execPath = path.join(execPath, 'resources');
execPath = path.join(execPath, 'app');

// select folder for current platform
switch (platform) {
	case 'darwin':
		execPath = path.join(execPath, 'svx.bin');
		break;
	case 'linux':
		execPath = path.join(execPath, 'svx.bin');
		break;
	case 'win32':
		execPath = path.join(execPath, 'svx.exe');
		isWin = true;
		break;
	default:
		global.console.log("unsupported platform: " + platform);
		break;
}
var svr = spawn(execPath.toString());
var win = gui.Window.get();
svr.stdout.on('data', function(data) {
    var msg = data.toString();
    //console.log(msg);
    win.window.location.href = msg;
});
win.on('loaded', function() {
	// the native onload event has just occurred
	this.show();
});
// Listen to the minimize event
win.on('minimize', function() {
	//console.log('Window is minimized');
	this.minimize(true);
});

win.on('close', function() {
	svr.kill();
    gui.App.quit();
});


