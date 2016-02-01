'use strict';

var app = require('app');
var BrowserWindow = require('browser-window');
var path = require('path');
var spawn = require('child_process').spawn;
//var arch = process.arch;
var platform = process.platform;

var isWin = false;


var execPath = path.dirname(process.execPath);
// some base path is appended
execPath = path.join(execPath, 'resources/app');

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
        //global.console.log("unsupported platform: " + platform);
        break;
}

var svr = null;

var mainWindow = null;

app.on('ready', function() {
    mainWindow = new BrowserWindow({
        title: "gDrv",
        frame: true,
        resizable: true,
        height: 600,
        width: 800
    });

    mainWindow.setMenu(null);
    //svr = spawn(tmpPath);
    svr = spawn(execPath.toString());
    svr.stdout.on('data', function(data) {
        var msg = data.toString();
        //console.log(msg);
        if(msg.match(/^http.*$/)) {
            mainWindow.loadURL(msg);
        }else{
            console.log("Master Thread Says: "+ msg);
        }
        // @TODO: Needto check if output is http then set window location else do something else

    });
    //mainWindow.loadURL('file://' + __dirname + '/app/index.html');
});

app.on('window-all-closed', function(){
    svr.kill();
    app.quit();
});


