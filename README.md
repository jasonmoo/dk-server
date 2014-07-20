
	    ._____
	  __| _/  | __   ______ ______________  __ ___________
	 / __ ||  |/ /  /  ___// __ \_  __ \  \/ // __ \_  __ \
	/ /_/ ||    <   \___ \\  ___/|  | \/\   /\  ___/|  | \/
	\____ ||__|_ \ /____  >\___  >__|    \_/  \___  >__|
	     \/     \/      \/     \/                 \/

	                       v 1,0


trend/cardinality tracking server

##HTTP Usage:
	// increment a key
	// g (group name)
	// k (key name)
	// v (optional; increment amount, defaults to 1)
	http://localhost?g=users&k=jason&v=1
	http://localhost?g=plays&k=catvideo&v=1

	// get the top n results for a group
	// g (group name)
	// n (count to return; capped at 200)
	http://localhost/top?g=users&n=100

	// Sample output:
	{
		"running": true,
		"timestamp": "2014-07-20 14:56:07.950249659 -0700 PDT",
		"render_time": "32.07us",
		"decay_rate": 0.05,
		"decay_floor": 1,
		"result_set": {
			"users": {
				"table_size": 1,
				"cardinality": {
					"percent": 0.2,
					"duration": "1m0s",
					"uniques": 1,
					"total": 5
				},
				"result_count": 1,
				"results": [
					{
						"name": "jason",
						"score": 4.494319519969077
					}
				]
			}
		}
	}

##Websocket Usage:
websockets are available for consuming reports regularly

	// console test:

	// open a connection
    var sock = new WebSocket("ws://localhost/sub");

    // simple bind for reporting on all events
	sock.onclose = sock.onerror = sock.onmessage = sock.onopen = console.log.bind(console);

	// subscribe to log(s)
	sock.send(JSON.stringify(["users","plays"]))

	// clear all subscriptions
	sock.send("[]")


