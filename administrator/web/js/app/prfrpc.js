
	prof = window.prof || {};
	
	prof.RPCEventHandlers = {};
	
	prof.handleSRPC = function(msg) {
		
		var response = JSON.parse(msg);
		
		if(response.error != null){
			console.log(response.error)
		}
		
		var result = response.result;
		
        if(typeof prof.RPCEventHandlers[response.id] == "undefined"){
        	console.log("Undefined binding: " + response.id);
        	return;
        } else {
        	// Pass to event handler.
        	prof.RPCEventHandlers[response.id](result);
        }
        
    }
    
	prof.MemStats = function(){
		var method = "ProfilerSRPC.MemStats";
		var params = {};
		postRPC(method,params);
	}
		
	window.prof = prof;