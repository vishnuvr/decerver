event = function(){
	
	var Event = {
		"Name" : "",
		"Source" : "",
		"Target": "",
		"Resource" : null,
		"Timestamp" : 0,
	};
	return evt;
}

var EventManager = {
	
	"callbacks" : {},
	
	"registerSubscriber" : function(source, subscriber){
		if(typeof callbackFn !== "function"){
			throw new Error("Callback '" + callbackF)
		}
		this.callbacks[source] = callbackFn
	},
};