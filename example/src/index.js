require('http').createServer(function (req, res) {
  res.end('Hello world!')
  console.log('Request received...')
}).listen(8080)

console.log('server started...')
