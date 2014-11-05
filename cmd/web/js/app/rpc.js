function postRPC(method, params) {
	var req = {
		"method" : method,
		"params" : params,
		"timestamp" : new Date().getTime()
	}
	var sfreq = JSON.stringify(req);
	console.log(sfreq);
	conn.send(sfreq);
}
