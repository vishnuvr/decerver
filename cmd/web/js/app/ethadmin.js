	var conn = null;

	var hasStorage = false;
	// Block chain table
	var bcTable;
	var bcTableApi;
	
	// Contract table
	var ctTable;
	var ctTableApi;
	
	// User table
	var usrTable;
	var usrTableApi;
	
	// For loading screen.
	var $modal;
	var bar;
	
	$('#LoadBCModal').modal({
		  backdrop: 'static',
		  show: false
	});
	
	$('#TxPool').enscroll({
	    showOnHover: true,
	    verticalTrackClass: 'trackV',
	    verticalHandleClass: 'handleV',
	    easingDuration : 100
	});
	
	$('#AccountCode').enscroll({
	    showOnHover: true,
	    verticalTrackClass: 'trackV',
	    verticalHandleClass: 'handleV',
	    easingDuration : 100
	});
	
	$('#AccountStorage').enscroll({
	    showOnHover: true,
	    verticalTrackClass: 'trackV',
	    verticalHandleClass: 'handleV',
	    easingDuration : 100
	});
	
	$('#BlockTransactions').enscroll({
	    showOnHover: true,
	    verticalTrackClass: 'trackV',
	    verticalHandleClass: 'handleV',
	    easingDuration : 100
	});
	
	$('#BlockUncles').enscroll({
	    showOnHover: true,
	    verticalTrackClass: 'trackV',
	    verticalHandleClass: 'handleV',
	    easingDuration : 100
	});
	
	$('#TerminalWidget').terminal(function(command, term) {
		if (command !== '') {
			try {
				var result = window.eval(command);
				if (result !== undefined) {
					term.echo(new String(result));
				}
			} catch (e) {
				term.error(new String(e));
			}
		} else {
			term.echo('');
		}
	}, {
		greetings : 'Ethereum console',
		name : 'eth_console',
		height : 500,
		prompt : '> '
	});
	
	var term = $.terminal.active();
	
	$("#BlockDataBody").hide();
	$("#EditorBody").hide();
	$("#AccountDataBody").hide();
	
	$("#menu-toggle").click(function(e) {
	    e.preventDefault();
	    
	    $("#wrapper").toggleClass("active");
	    
	    if($('#main_icon').hasClass("fa-arrow-circle-o-left")){
	    	$('#main_icon').removeClass("fa-arrow-circle-o-left");
	    	$('#main_icon').addClass("fa-arrow-circle-o-right");
	    } else {
	    	$('#main_icon').removeClass("fa-arrow-circle-o-right");
	    	$('#main_icon').addClass("fa-arrow-circle-o-left");
	    }
	});
	
	$(".toggle-button").click(function(e) {
		var tb = $(this);
	    e.preventDefault();
	    if(tb.hasClass("fa-minus")){
	    	tb.removeClass("fa-minus");
	    	tb.addClass("fa-plus");
	    } else {
	    	tb.removeClass("fa-plus");
	    	tb.addClass("fa-minus");
	    }
	    
	    var body = $(this).parent().parent().parent().children('.panel-body');
	    
	    if (body.hasClass('collapsed')) {
	        // expand the panel
	    	body.slideDown();
	    	body.removeClass('collapsed');
	        if(body.attr("id") === "EditorBody" && !DBstarted){
	        	DBstarted = true;
	        	renderDB(important,0);
	        }
	    }
	    else {
	        // collapse the panel
	    	body.slideUp();
	    	body.addClass('collapsed');
	    }
	});
	
	$(window).unload(function() {
		if(conn != null && conn.readyState < 2){
			conn.close();
		}
	});
	
	$( document ).ready(function() {
				
		prepTransactionPanel();
		
		bar = $('.progress-bar');
				
		bcTable = $('#blockchain-table').dataTable({
			"data" : [],
			"order" : [ [ 0, "desc" ] ],
			"deferRender" : true,
			"columnDefs": [
			               {
			                   "render": function ( data, type, row ) {
			                       return data.substring(0,10) + "..." ;
			                   },
			                   
			                   "targets": 2
			               }]
		});
		
		bcTableApi = $( "#blockchain-table" ).dataTable().api();
		
		ctTable = $('#contract-table').dataTable({
			"data" : [],
			"order" : [ [ 0, "asc" ] ],
			"deferRender" : true,
			"columnDefs": [
	               {
	                   "render": function ( data, type, row ) {
	                       return data.substring(0,10) + "..." ;
	                   },
	                   
	                   "targets": 0
	               },{
	                   "render": function ( data, type, row ) {
	                	   var disp = eth.convertWeiToBest(data);
	                	   if(disp[0] == "0"){
	                		   return "0"
	                	   }
	                	   if(disp[0].length > 6){
	                		   disp[0] = disp[0].substring(0,6) + "...";
	                	   }
	                       return disp[0] + " " + disp[1];
	                   },
	                   
	                   "targets": 2
	               }]
		});
		
		ctTableApi = $( "#contract-table" ).dataTable().api();
		
		usrTable = $('#user-table').dataTable({
			"data" : [],
			"order" : [ [ 0, "asc" ] ],
			"deferRender" : true,
			"columnDefs": [
			               {
			                   "render": function ( data, type, row ) {
			                       return data.substring(0,10) + "..." ;
			                   },
			                   
			                   "targets": 0
			               },{
			                   "render": function ( data, type, row ) {
			                	   var disp = eth.convertWeiToBest(data);
			                	   
			                	   if(disp[0] == "0"){
			                		   return "0"
			                	   }
			                	   
			                	   if(disp[0].length > 6){
			                		   disp[0] = disp[0].substring(0,6) + "...";
			                	   }
			                	   
			                       return disp[0] + " " + disp[1];
			                   },
			                   
			                   "targets": 2
			               }]
		});
		
		usrTableApi = $( "#user-table" ).dataTable().api();
				
		$('#blockchain-table tbody').on( 'click', 'tr', function () {
			var hash = bcTable.fnGetData( this )[2];
			eth.BlockByHash(hash)
		} );
		
		$('#contract-table tbody').on( 'click', 'tr', function () {
			var addr = ctTable.fnGetData( this )[0];
			eth.Account(addr);
		} );
		
		$('#user-table tbody').on( 'click', 'tr', function () {
			var addr = usrTable.fnGetData( this )[0];
			eth.Account(addr);
		} );
				
		initRPCEventHandlers();
		
		if (window["WebSocket"]) {
	        conn = new WebSocket("ws://localhost:3000/wsapi");
	        
	        conn.onopen = function () {
	        	$modal = $('#LoadBCModal');
	    		$modal.modal('show');
	        	console.log("WebSockets connection opened.");
	        	eth.Init();
	        }
	        
	        conn.onclose = function(evt) {
	            console.log("WebSockets connection closed.");
	        }
	        
	        conn.onmessage = function(evt) {	
	        	eth.handleSRPC(evt.data);
	        }
	    } else {
	    	alert("Your browser does not support WebSockets.\nMonk admin needs websockets in order to work.");
	    }
		
	});
	
	function addBlockToChain(block){
		bcTable.fnAddData(block);
		bcTable.fnDraw();
	}
	
	function addBlocksToChain(blocks){
		for(var i = 0; i < blocks.length; i++){
			bcTable.fnAddData(blocks[i]);
		}
		
		bcTable.fnDraw();
	}
	
	
	
	function prepTransactionPanel(){
		
		$('#InputValue').val("0");
		$('#InputValueDenom').html("ether");
		var sList = "";
		for(var i = 0; i < eth.denomArr.length; i++){
			sList += '<li><a href="#" >' + eth.denomArr[i] + '</a></li>';
		}
		$('#InputValueDenomDropDown').html(sList);
		
		$('#InputGas').val("0");
		//var mingas = eth.MinGascost().SVal;
		//mingas = eth.convertWeiToBest(mingas);
		$('#InputGasPrice').val(10);
		$('#InputGasDenom').html("szabo");
		sList = "";
		for(var i = 0; i < eth.denomArr.length; i++){
			sList += '<li><a href="#" >' + eth.denomArr[i] + '</a></li>';
		}
		
		$('#InputGasDenomDropDown').html(sList);
	}
	
	$(function() {
		// Triggers		
	
		$("#ToggleMining").click(function(event) {
			var ml = $("#MiningLight");
			var tmi = $("#ToggleMiningIcon");
			
			if(ml.hasClass("led-green")){
				ml.removeClass("led-green");
				ml.addClass("led-red");
				tmi.removeClass("fa-stop");
				tmi.addClass("fa-play");
				
				eth.StopMining(false);
			} else {
				ml.removeClass("led-red");
				ml.addClass("led-green");
				tmi.removeClass("fa-play");
				tmi.addClass("fa-stop");
				eth.StartMining(false);
			}
		});
		
		$('#InputValueDenomDropDown li a').click(function(event){
			event.preventDefault();
			$('#InputValueDenom').html($(this).text());
			$(this).parent().parent().dropdown('toggle');
			return false;
		});
		
		$('#InputGasDenomDropDown li a').click(function(event){
			event.preventDefault();
			$('#InputGasDenom').html($(this).text());
			$(this).parent().parent().dropdown('toggle');
			return false;
		});
		
		
		$('#TxSendButton').click(function(event){
			event.preventDefault();
			// Receiver address
			var recipient = $('#InputRecipient').val();
			if(recipient === null){
				recipient = "";
			}
			
			// Value of transaction
			var value = $('#InputValue').val();
			var valDenom = $('#InputValueDenom').html().trim();
			
			try{
				value = (bigInt(value).times(bigInt(eth.denom[valDenom]))).toHex();
			} catch (err) {
				window.alert("Value or denomination is wrong.");
				return false;
			}
			
			// Gas
			var gas = $('#InputGas').val();
			gas = bigInt(gas).toHex();
			var gasCost = $('#InputGasCost').val();
			var gasCostDenom = $('#InputGasDenom').html().trim();
			
			try{
				gasCost = (bigInt(gasCost).times(eth.denom[gasCostDenom])).toHex();
				console.log("TotalGasCost: " + gasCost);
			} catch (err) {
				window.alert("Gas or gascost is wrong.");
				return false;
			}
			
			var data = $('#DataTextArea').val();
			
			var tx = eth.Transact(recipient,value,gas,gasCost,data);
			return false;
		});
		
		$('#TxClearButton').click(function(event){
			event.preventDefault();
			prepTransactionPanel();
			return false;
		});
	});
	
