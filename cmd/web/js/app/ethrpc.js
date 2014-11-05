
	eth = window.eth || {};

	var RPC_SERVICE = "MonkWsAPI.";
	
	// TODO Move away from these.
	var ERR_NO_SUCH_BLOCK = "NO SUCH BLOCK";
	var ERR_NO_SUCH_TX = "NO SUCH TX";
	var ERR_NO_SUCH_ADDRESS = "NO SUCH ADDRESS";
	var ERR_STATE_NO_STORAGE = "STATE NO STORAGE";
	var ERR_MALFORMED_ADDRESS = "MALFORMED ADDRESS";
	var ERR_MALFORMED_BLOCK_HASH = "MALFORMED BLOCK HASH";
	var ERR_MALFORMED_TX_HASH = "MALFORMED TX HASH";
	
	var ZEROADDR = "0000000000000000000000000000000000000000";
	
	// Account flags
	var ACCOUNT_MODIFIED = 0;
	var ACCOUNT_CREATED = 1;
	var ACCOUNT_DELETED = 2;
	
	eth.latestBlockNr = 0;
	eth.numAccounts = 0;
			
	eth.contracts = {};
	eth.users = {};
	
	eth.RPCEventHandlers = {};
	
	eth.handleSRPC = function(msg) {
		
		var response = JSON.parse(msg);
		
		if(response.error != null){
			console.log(response.error)
		}
		
		var result = response.result;
		
        if(typeof eth.RPCEventHandlers[response.id] == "undefined"){
        	console.log("Undefined binding: " + response.id);
        	return;
        } else {
        	console.log("Received: " + response.id);
        	// Pass to event handler.
        	eth.RPCEventHandlers[response.id](result);
        }
        
    }
	
    eth.Init = function(){
    	eth.MyAddress()
    	eth.MyBalance()
    	eth.WorldState();
    }
    
	eth.LastBlockNumber = function(){
		var method = RPC_SERVICE + "LastBlockNumber";
		var params = {};
		postRPC(method,params);
	}
	
	eth.WorldState = function(){
		var method = RPC_SERVICE + "WorldState";
		var params = {};
		postRPC(method,params);
	}

	/**
	 * Get the balance of the active account.
	 */
	eth.MyBalance = function(){
		var method = RPC_SERVICE + "MyBalance";
		var params = {};
		postRPC(method,params);
	}
	
	eth.MyAddress = function(){
		var method = RPC_SERVICE + "MyAddress";
		var params = {};
		postRPC(method,params);
	}
	
	/**
	 * Signals the Ethereum client to start mining.
	 * Returns true if mining was successfully started by running this command.
	 */
	eth.StartMining = function(forced){
		var method = RPC_SERVICE + "StartMining";
		var params = {};
		postRPC(method,params);
	}
	
	/**
	 * Signals the Ethereum client to start mining.
	 * Returns true if mining was successfully started by running this command.
	 */
	eth.StopMining = function(forced){
		var method = RPC_SERVICE + "StopMining";
		var params = {};
		postRPC(method,params);
	}
	
	/**
	 * Gets a block from its hash.
	 * Returns a block object. 
	 */
	eth.BlockByHash = function(hash){
		var method = RPC_SERVICE + "BlockByHash";
		var params = {"SVal" : hash};
		postRPC(method,params);
	}
	
	eth.Transact = function(recipient,value,gas,gascost,data){
		var method = RPC_SERVICE + "Transact";
		var params = {	"Recipient" : recipient,
						"Value"		: value,
						"Gas"		: gas,
						"GasCost"	: gascost,
						"Data"		: data,
					};
		postRPC(method,params);
	}
	
	
	eth.MinGascost = function(){
		var method = RPC_SERVICE + "MinGascost";
		var params = {};
		postRPC(method,params);
	}
	
	eth.Account = function(addr){
		var method = RPC_SERVICE + "Account";
		var params = {"SVal" : addr};
		postRPC(method,params);
	}
	
	window.eth = eth;
