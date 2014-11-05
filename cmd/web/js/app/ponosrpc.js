
	ponos = window.ponos || {};
	
	ponos.RPCEventHandlers = {};
	
	ponos.handleSRPC = function(msg) {
		
		var response = JSON.parse(msg);
		
		if(response.error != null){
			window.alert(response.error)
		}
		
		var result = response.result;
		
        if(typeof ponos.RPCEventHandlers[response.id] == "undefined"){
        	console.log("Undefined binding: " + response.id);
        	return;
        } else {
        	// Pass to event handler.
        	ponos.RPCEventHandlers[response.id](result);
        }
        
    }
    
	ponos.GetTree = function(){
		var method = "PonosSRPC.GetTree";
		var params = {};
		postRPC(method,params);
	}
		
	window.ponos = ponos;