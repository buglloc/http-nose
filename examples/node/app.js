var express = require('express');
var bodyParser = require('body-parser');

var app = express();
app.use(bodyParser.raw());

app.use(function(req, res) {
	let picked = req.header('x-pick');
	if (! picked) {
        picked = []
	} else if (typeof picked === 'string') {
    	picked = [picked];
    }

	let data = {
		'Method': req.method,
		'RequestURI': req.originalUrl,
		'Path': req.path,
		'Args': req.query,
		'Headers': Object.keys(req.headers).map((key) => {
			return {'Name': key, 'Value': req.headers[key]}
		}),
		'Body': req.body,
		'Host': req.hostname
	}

	res.setHeader('Content-Type', 'application/json');
	res.send(JSON.stringify(data, null, 2));
});

app.listen(3000, function () {
  console.log('Example app listening on port 3000!');
});
