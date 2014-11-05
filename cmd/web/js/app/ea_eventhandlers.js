
var myTxs = [];
var postTxQueue = {};

initRPCEventHandlers = function(){
	// Done at startup to get the active address.
	eth.RPCEventHandlers["MyBalance"] = updateMyBalance;
	// Done at startup to get the user address.
	eth.RPCEventHandlers["MyAddress"] = updateMyAddress;
	// Not used atm.
	eth.RPCEventHandlers["LastBlockNumber"] = lastBlockNumber;
	// Done when clicking a block in the block table.
	eth.RPCEventHandlers["BlockByHash"] = updateBlockView;
	// Done when clicking an account in the contracts or user tables.
	eth.RPCEventHandlers["Account"] = updateAccount;
	// Reactor event.
	eth.RPCEventHandlers["BlockAdded"] = blockAdded;
	// Mining on.
	eth.RPCEventHandlers["StartMining"] = consumeEvent;
	// Mining off.
	eth.RPCEventHandlers["StopMining"] = consumeEvent;
	
	// Transactions
	eth.RPCEventHandlers["Transact"] = handleTx;
	eth.RPCEventHandlers["TxPre"] = handleTxPre;
	eth.RPCEventHandlers["TxPreFail"] = handleTxPreFail;
	eth.RPCEventHandlers["TxPost"] = handleTxPost;
	eth.RPCEventHandlers["TxPostFail"] = handleTxPostFail;
	
	// These events happen when doing an eth.WorldState() on startup, 
	// and happens in this order:
	//
	// NumBlocks - sets eth.lastBlockNumber and prepares the system for
	// receiving that amount of blocks. The server then starts sending
	// blocks.
	//
	// Blocks - one event per block.
	// 
	// NumAccounts - sends the total number of accounts.
	//
	// Accounts - one per account
	//
	// WorldStateDone - signals that the webclient has received the
	// entire world state.
	eth.RPCEventHandlers["NumBlocks"] = numBlocks;
	eth.RPCEventHandlers["Blocks"] = blocks;
	eth.RPCEventHandlers["NumAccounts"] = numAccounts;
	eth.RPCEventHandlers["Accounts"] = accounts;
	eth.RPCEventHandlers["WorldStateDone"] = worldStateDone;
	
	eth.RPCEventHandlers["Log"] = printLogMessage;
}

numBlocks = function(result){
	eth.lastBlockNumber = result.IVal;
}

blocks = function(result){
	var num = result.Number;
	bcTable.fnAddData([result.Number,result.Transactions,result.Hash]);
	if(eth.lastBlockNumber != 0){
		var val = (100*num)/eth.lastBlockNumber;
		bar.css('width', val+'%').attr('aria-valuenow', val);
	}
	if(num == eth.lastBlockNumber){
		bcTable.fnDraw();
		$modal.modal('hide');
	}
};

numAccounts = function(result){
	eth.numAccounts = result.IVal;
}

accounts = function(accountMini){
	if(accountMini.Contract){
		eth.contracts[accountMini.Address] = accountMini
		ctTable.fnAddData([accountMini.Address,accountMini.Nonce,accountMini.Value]);
	} else {
		eth.users[accountMini.Address] = accountMini
		usrTable.fnAddData([accountMini.Address,accountMini.Nonce,accountMini.Value]);
	}
}

worldStateDone = function(){
	ctTable.fnDraw();
	usrTable.fnDraw();
	$modal.modal('hide');
	console.log("World state updated. Last block: " + eth.lastBlockNumber + ". Accounts: " + eth.numAccounts)
}

lastBlockNumber = function(result){
	eth.lastBlockNumber = result.IVal;
}

updateMyBalance = function(balance){
	var myEther = eth.convertWei("0x" + balance.SVal, "ether", 2);
	$('#MyEther').html(myEther);
}

updateMyAddress = function(address){
	$('#MyAddress').html(address.SVal);
}

updateBlockView = function(block){
	if($('#BlockDataBody').hasClass('collapsed')){
		$('#BlockDataToggle').click();
	}
	
	$( "#BlockNumber" ).html( block.Number );
	$( "#BlockTransNr" ).html( block.Transactions.length );
	$( "#BlockTime" ).html( eth.timestampToDate(block.Time) );
	var bmg = eth.convertWeiToBest(block.MinGasPrice)
	$( "#BlockMinGas" ).html( bmg[0] + " " + bmg[1] );
	$( "#BlockGasLimit" ).html( block.GasLimit );
	$( "#BlockGasUsed" ).html( block.GasUsed );
	$( "#BlockDifficulty" ).html( block.Difficulty );
	
	// Extras
	$( "#BlockHash" ).html( block.Hash );
	$( "#BlockPrevHash" ).html( block.PrevHash );
	$( "#BlockNonce" ).html( block.Nonce );
	$( "#BlockCoinbase" ).html( block.Coinbase );
	
	var txSha;
	if(block.TxSha !== ""){
		txSha = block.TxSha;
	} else {
		txSha = "-";
	}
	$( "#BlockTxSha" ).html( txSha );
	$( "#BlockUncleSha" ).html( block.UncleSha );
			
	var liString = "";
	// Now we add transactions
	for(var i = 0; i < block.Transactions.length; i++){
		var tx = block.Transactions[i];
		var from = "From:"
		var to = "To:";
		var extra = "Normal Transaction"
		if(tx.ContractCreation){
			from = "Creator:"
			to = "Address:";
			extra = "Contract Creation"
		}
		
		
		liString += '<li class="list-group-item">';
		liString += getTxHeader(i + 1,extra);
		liString += get2_10keyValueDiv(from,tx.Sender);
		liString += get2_10keyValueDiv(to,tx.Recipient);
		liString += get2_10keyValueDiv("Hash:", tx.Hash);
		liString += get2_10keyValueDiv("Nonce:", tx.Nonce);
		liString += get2_10keyValueDiv("Gas:", tx.Gas);
		liString += get2_10keyValueDiv("GasCost:", tx.GasCost);
		liString += '</li>'; 
	}
	
	$( "#BlockTransactions" ).html( liString );
			
	// Uncles
	liString = "";
	// Now we add transactions
	for(var i = 0; i < block.Uncles.length; i++){
		liString += '<li class="list-group-item">' + block.Uncles[i].Hash + '</li>'; 
	}
	$( "#BlockUncles" ).html( liString );
}

updateAccount = function(account){
	
	if($('#AccountDataBody').hasClass('collapsed')){
		$('#AccountDataToggle').click();
	}
	
	$('#AccountAddress').html( account.Address );
	$('#AccountNonce').html( account.Nonce );
	$('#AccountValue').html( account.Value );
	if(account.Code !== null && account.Code !== ""){
		var str = "";
		
		for (var i = 0; i < account.Code.length; i += 2)
		{
		    str += "0x" + account.Code.substring(i,i + 2) + " ";
		}
		$('#AccountCode').html( str );
	} else {
		$('#AccountCode').html( "" );
	}
	
	if(account.Storage.length > 0){
		var liString = "<samp>";
		// Now we add transactions
		for(var i = 0; i < account.Storage.length; i += 2){
			var st = account.Storage[i];
			var key = "Addr: " + eth.stringForDisplay(st);
			
			st = account.Storage[i + 1];
			var val = "<b>Val:</b> " + eth.stringForDisplay(st);
			
			liString += '<li class="list-group-item">';
			liString += get12_2r_keyValueDiv(key,val);
			liString += '</li>'; 
		}
		liString += "</samp>";
		$( "#AccountStorage" ).html( liString );
	} else {
		$('#AccountStorage').html( "" );
	}
}

blockAdded = function(result){
	// TODO error handling here, in case we're out of sync.
	var num = result.Number;
	bcTable.fnAddData([result.Number,result.Transactions,result.Hash]);
	eth.lastBlockNumber++;
	bcTable.fnDraw();
	
	var ctRedraw = false, usrRedraw = false;
	// Now update accounts
	for (var i = 0; i < result.AccountsAffected.length; i++) {
		var acc = result.AccountsAffected[i];
		console.log(acc);
		// Check if it's a new contract.
		if(acc.Flag & 1 == 1){
			// TODO Remove after bug testing.
			if(typeof eth.contracts[acc.Address] != "undefined"){
				console.log("Buggy already-existing contract with create flag");
				break;
			}
			eth.contracts[acc.Address] = 1;
			ctTable.fnAddData([acc.Address,acc.Nonce,acc.Value]);
			ctRedraw = true;
			eth.numAccounts += 1;
		} else if(typeof eth.contracts[acc.Address] != "undefined") {
			// This is a contract. Is it being deleted?
			if(acc.Flag & 2 == 2){
				eth.contracts[acc.Address] = null;
				var index = usrTableApi.column( 0 ).data().indexOf( acc.Address );
				ctTableApi.remove(index);
				eth.numAccounts -= 1;
			} else {
				// It's being modified.
				var index = ctTableApi.column( 0 ).data().indexOf( acc.Address );
				ctTable.fnUpdate([acc.Address,acc.Nonce,acc.Value],index);
			}
			ctRedraw = true;
		} else if(typeof eth.users[acc.Address] != "undefined") {
			// This is a user account being modified. Since users cannot be deleted, 
			// and it already exists in the system, it must be a modification.
			var index = usrTableApi.column( 0 ).data().indexOf( acc.Address );
			console.log(index)
			usrTable.fnUpdate([acc.Address,acc.Nonce,acc.Value],index);
			usrRedraw = true;
		} else {
			// This is a brand new user.
			eth.users[acc.Address] = 1;
			usrTable.fnAddData([acc.Address,acc.Nonce,acc.Value]);
			usrRedraw = true;
			eth.numAccounts += 1;
		}
	}
	if(ctRedraw){
		ctTable.fnDraw();
	}
	if(usrRedraw){
		usrTable.fnDraw();
	}
	eth.MyBalance()
}

handleTx = function(result){
	// We get a recipe back. If this step fails it means that the compilation
	// of code fails, for example, or it could be some weird Ethereum error.
	if (result.Error == ""){
		myTxs.push[result.Hash];
	} else {
		window.alert("TX Error: " + result.Error);
	}
}

handleTxPre = function(result){
	// Add this transaction to the transaction window. Its status now is "queued".
	$('#TxPool').prepend(getTxPoolLi(result));
}

handleTxPreFail = function(result){
	var theIdx = -1;
	for(var i = 0; i < myTxs.length; i++){
		var hash = myTxs[i];
		if(hash === result.Hash){
			theIdx = i;
			window.alert("TX Error: " + result.Error);
			break;
		}
	}
	if(theIdx >= 0){
		myHashes.splice(theIdx);
	}
}

handleTxPost = function(result){
	var listEntry = document.getElementById('TX_' + result.Hash);

	console.log("PostTx hash: " + result.Hash);
	if(listEntry === null){
		console.log("TxPost: Transaction not stored in list: " + result.Hash);
	} else {
		var status = document.getElementById('status_TX_' + result.Hash);
		// TODO Second check not needed.
		if (typeof postTxQueue[result.Hash] === "undefined" || postTxQueue[result.Hash] === null) {		
			listEntry.style.backgroundColor = "#FFA";
			status.innerHTML = "Status: PENDING"
			// Just add something.
			postTxQueue[result.Hash] = "waiting";
			
		}
		else {
			listEntry.style.backgroundColor = "#AFA";
			status.innerHTML = "Status: SUCCESS"
			postTxQueue[result.Hash] = null;
			var theIdx = -1;
			for(var i = 0; i < myTxs.length; i++){
				var hash = myTxs[i];
				if(hash === result.Hash){
					theIdx = i;
					break;
				}
			}
			if(theIdx >= 0){
				myTxs.splice(theIdx);
			}
		}
	}
}

handleTxPostFail = function(result){
	var listEntry = document.getElementById('TX_' + result.Hash);
	console.log(listEntry)
	if(listEntry === null){
		console.log("TxPost: Transaction not stored in list: " + result.Hash);
	} else {
		listEntry.style.backgroundColor = "#faa";
		document.getElementById('status_TX_' + result.Hash).innerHTML = "Status: REJECTED" + "<br>Error: " + result.Error;
	}
	var theIdx = -1;
	for(var i = 0; i < myTxs.length; i++){
		var hash = myTxs[i];
		if(hash === result.Hash){
			theIdx = i;
			window.alert("TX Error: " + result.Error);
			break;
		}
	}
	if(theIdx >= 0){
		myTxs.splice(theIdx);
	}
}

printLogMessage = function(result){
	term.echo(result.SVal)
}

// Used if response does not need handling.
consumeEvent = function(result){
	console.log("SRPC message not handled.")
}

// Helpers
function getTxHeader(num, extra){
	return '<div class="row"></div><div class="col-sm-12"><h4>#' + num.toString() + ' (' + extra + ')</h4></div></div>';
}

function get2_10keyValueDiv(key, value){
	return '<div class="row"><div class="col-sm-2"><b>' + key.toString() + '</b></div><div class="col-sm-10">' + value.toString() + '</div></div>';
}

function get12_2r_keyValueDiv(key, value){
	return '<div class="row"><b>' + key.toString() + '</b></div><div class="row">' + value.toString() + '</div>';
}

function getTxPoolLi(txData){
	var cc = txData.ContractCreation ? "+" : "*";
	var txFromTo = txData.Sender.substring(0,15) + '... (' + cc + ')-> ' + txData.Recipient.substring(0,15) + '...';
	var value = eth.convertWeiToBest(txData.Value);
	if(value[0].length > 15){
		value[0] = value[0].substring(0,15) + "...";
	}
	var id ='"TX_' + txData.Hash + '"';
	var li = '<li id=' + id + '><h4>Transaction</h4>Info: ' + txFromTo + '<br>Value: ' + value[0] + ' ' + value[1] + '<br><div id="status_TX_' + txData.Hash + '">Status: QUEUED</div>' + '</li>';
	console.log("Tx list entry: " + li);
	console.log("Writing transaction to id: " + id);
	return li;
}
