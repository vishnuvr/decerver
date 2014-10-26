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
	
	// Code mirror object
	var editor;
	
	var frSupport;
	
	// Code mirror theme selector
	var input = document.getElementById("ThemeSelector");
	
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
	$("#TxPoolBody").hide();
	
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
					
		// Check if localstorage exists.
		if(typeof(window.localStorage) !== "undefined") {
		    hasStorage = true;
		}
		
		frSupport = window.File && window.FileReader && window.FileList;
		
		prepTransactionPanel();
		editor = CodeMirror.fromTextArea(document.getElementById("DataTextArea"),{
	        matchBrackets: true,
	        indentUnit: 4,
	        tabSize: 4,
	        indentWithTabs: true,
	        lineNumbers: true,
	        autoCloseBrackets: true,
	        matchBrackets: true,
	        styleActiveLine: true,
	        mode: "text/x-go",
			extraKeys : {
				"F11" : function(cm) {
					cm.setOption("fullScreen", !cm.getOption("fullScreen"));
				},
				"Esc" : function(cm) {
					if (cm.getOption("fullScreen"))
						cm.setOption("fullScreen", false);
				},
				"Ctrl-S" : function(){
					editor.save();
				} 
			}
	    });
		
		editor.save = function(){
			var str = editor.getValue();
			if(str === null || str.length === ""){
				return;
			}
			saveTextAs(str, "etheris.mu");
		}
		
		editor.open = function(){
			
		}
		
		var choice = "rubyblue";
		// Load settings.
		if(hasStorage && localStorage.editorTheme){
			choice = localStorage.editorTheme;
		} else if (hasStorage){
			// Init
			localStorage.editorTheme = choice;
		}
		$('#ThemeSelectorText').html(choice);
		editor.setOption("theme", choice);
		
		// Start downloading the world state.
		
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
	        conn = new WebSocket("wss://cledus.erisindustries.com/srpc");
	        
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
	    	alert("Your browser does not support WebSockets.\nEthereum admin needs websockets in order to work.");
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
		
		sList = "";
		
		sList += '<li><a href="#" >default</a></li>';
		sList += '<li><a href="#" >3024-day</a></li>';
		sList += '<li><a href="#" >3024-night</a></li>';
		sList += '<li><a href="#" >ambiance</a></li>';
		sList += '<li><a href="#" >base16-dark</a></li>';
		sList += '<li><a href="#" >base16-light</a></li>';
		sList += '<li><a href="#" >blackboard</a></li>';
		sList += '<li><a href="#" >cobalt</a></li>';
		sList += '<li><a href="#" >eclipse</a></li>';
		sList += '<li><a href="#" >elegant</a></li>';
		sList += '<li><a href="#" >erlang-dark</a></li>';
		sList += '<li><a href="#" >lesser-dark</a></li>';
		sList += '<li><a href="#" >mbo</a></li>';
		sList += '<li><a href="#" >mdn-like</a></li>';
		sList += '<li><a href="#" >midnight</a></li>';
		sList += '<li><a href="#" >monokai</a></li>';
		sList += '<li><a href="#" >neat</a></li>';
		sList += '<li><a href="#" >neo</a></li>';
		sList += '<li><a href="#" >night</a></li>';
		sList += '<li><a href="#" >paraiso-dark</a></li>';
		sList += '<li><a href="#" >paraiso-light</a></li>';
		sList += '<li><a href="#" >pastel-on-dark</a></li>';
		sList += '<li><a href="#" >rubyblue</a></li>';
		sList += '<li><a href="#" >solarized dark</a></li>';
		sList += '<li><a href="#" >solarized light</a></li>';
		sList += '<li><a href="#" >the-matrix</a></li>';
		sList += '<li><a href="#" >tomorrow-night-eighties</a></li>';
		sList += '<li><a href="#" >twilight</a></li>';
		sList += '<li><a href="#" >vibrant-ink</a></li>';
		sList += '<li><a href="#" >xq-dark</a></li>';
		sList += '<li><a href="#" >xq-light</a></li>';
		
		$('#ThemeSelectorDropDown').html(sList);
		
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
		
		$('#ThemeSelectorDropDown li a').click(function(event){
			event.preventDefault();
			$('#ThemeSelectorText').html($(this).text());
			$(this).parent().parent().dropdown('toggle');
			editor.setOption("theme", $(this).text());
			if (hasStorage){
				localStorage.editorTheme = $(this).text();
			}
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
			
			var data = editor.getValue();
			
			var tx = eth.Transact(recipient,value,gas,gasCost,data);
			return false;
		});
		
		$('#TxClearButton').click(function(event){
			event.preventDefault();
			prepTransactionPanel();
			return false;
		});
		
		$('#EditorUndoButton').click(function(event){
			event.preventDefault();
			editor.undo();
			return true;
		});
		
		$('#EditorRedoButton').click(function(event){
			event.preventDefault();
			editor.redo();
			return true;
		});
		
		$('#EditorSaveButton').click(function(event){
			event.preventDefault();
			editor.save();
			return false;
		});
		
		$('#EditorOpenButton').click(function(event){
			event.preventDefault();
			if (frSupport) {
				$('#FileOpenBox').click();
			} else {
			  alert('The File APIs are not fully supported in this browser. You must use Copy/Paste.');
			}
			return false;
		});
		
		$('#CompileButton').click(function(event){
			event.preventDefault();
			
			var mutanStr = editor.getValue();
			if(mutanStr === null || mutanStr.length === 0){
				$('#CompilerOutput').html("<samp><b>Nothing to compile.</b></samp>");
				return false;
			}
			
			eth.CompileMutan(mutanStr);
			
			return false;
		});
		
		document.getElementById('FileOpenBox').addEventListener('change', handleFileSelect, false);
	});
	
	function handleFileSelect(evt) {
		
		var files = evt.target.files; // FileList object
		var file = files[0];
		
		var reader = new FileReader();
	
		// Closure to capture the file information.
		reader.onload = (function(theFile) {
			return function(e) {
				editor.setValue(e.target.result);
			};
		})(file);
	
		// Read in the image file as a data URL.
		reader.readAsText(file);
	}
	
	var DBdone = false;
	var DBstarted = false;
	
	var important = [];
	important.push("                                       m$$@$@@$@@$@@@,");
	important.push("                                      n@@@%@#jC](x @@@@#@? ");                                         
	important.push("                                   x@@%@$@x          ;@@@@@ ");                                        
	important.push("                                 $@#@@@_ .             @@@@@@. ");
	important.push("                         -Q@$@@@$$@                      @@@$@@ ");
	important.push("                      ($@@@%$@@$@@?.                     -@#@@@ ");
	important.push("                     @@@#@@#@@@W1$@                       x@@@@+ ");
	important.push("                    >@@@@@@&$&@$ @@|                       @@#@n ");
	important.push("                    $@@@@@@@@@@z !@@i                      <@@$] ");
	important.push("                    @@@@@@@$@@@  ,$@f         x@@$@Y}       !@%r ");
	important.push("                    @@@@@@@$@    n#@]#@<    v@@@@@@@@@$?    ]@@c ");
	important.push("                    @@@@@@@B?i-_v#@@@#@    $@@@@#@@$$i@@#   1@@[ ");
	important.push("                    {@@@&@@@@$%@#@@f      v@@@@@@@$@@  @@@U j@$x ");
	important.push("                    @@@     J@@@@  Jc    j@@@@@@@$@#@  t$#@.C@$1 ");
	important.push("                   @$@L    x@%@@   ]@p   #@#$@@@@@@##  !@@@ #@@+ ");
	important.push("                   @@8     _@@$@@@@@$@<  #@@@@@$@@@B.   $@iY$@@ ");
	important.push("                 @@@@@@@@*(l     .       ~@@@@mI       %@  @$#@ ");
	important.push("                !@@@$@@@@@@@@@@@}l         @@@@@@@@@@@#h   @@@# ");
	important.push("               _@@@( .I!p$@@@@@@@@@$@kJ       1&@$#@@@x   )@@@8 ");
	important.push("               ##@$          (tB@@$@@##@#$@ti             B#@#J ");
	important.push("              z@#@_                    ;!;t###@@@U1      }@@@# ");
	important.push("             r@#@?                                       #@$@n ");
	important.push("             @@@@                                        @$@@ ");
	important.push("            >@@@L.                                      z@@@@ ");
	important.push("           ~$@@!                                       ]@@#@n ");
	important.push("           @$@#                                        ##@&Q ");
	important.push("          x@@@u                                        $@$@ ");
	important.push("         I@@#U                                        |@@@$ ");
	important.push("         (@@@.                          |@#          |@@$; ");
	important.push("         [@@@                           @@@   @@<    @@@@ ");
	important.push("         n@@@                          U@@]  @@@J    @&@x                               1aqU> ");      
	important.push("         j@@@                          @@$>  @@@|   (@@@                              J@$W#@@@t  ");  
	important.push("         J@@m                          @@}  z@@1   1@@@                             q@@$     @@M   ");
	important.push("         J@@x                          @$. {@@@.   @@@@                            $@@@      @@@.  ");
	important.push("         J@#Z                         v@$  @#@Z    @#@,                          Z@@$%       @@C   ");
	important.push("         J@@Y                         k$Z  @##    {@&@,                         }@@@@        @%    ");
	important.push("         J@@C                        f@]  >@@    i@@#n                         @@#@        i#@J    ");
	important.push("         J@@U                        @@  X@@@    t##@                        o@@@x         @@i     ");
	important.push("         J@@Y                        @@  @@@     z@@@  -*@#@@@@@@@@@J       Z$@@n        ~@@@      ");
	important.push("         z@@{                        @0 {@#~     f#@@@@@$@$@@@$@@@@@$@@    @@@@ .        @@C       ");
	important.push("         )#@$                       ;8 c@@&      j$$@l    ~jI      1w@@$@M#$@          @@@C        ");
	important.push("         !@@#                      +@! @@@       1%@#      @@@-       >%@@@          Z@#@>         ");
	important.push("          @@#                      ##  @$)         ,;      @#@/        @@]          !#@@*          ");
	important.push("          W@@                      @$  @&                  uUW8@@BQx>m$@z          #@$@L           ");
	important.push("          [@@)                   @@;    X@h            @@@@@@@%@#@@##$@@I        @@@@@j            ");
	important.push("           @@b                 #Z        U@]            !1/_,;_(rx##@@#@#@u     i###@U_            ");
	important.push("           ###]               @      [@t  .@v.                      .|%@#@@x       #@@@@L          ");
	important.push("           a@@@               #@@@@  @@@$ .@?                          {@@@@@         _i@$@        ");
	important.push("           <@@@!                 @@Lj@@.;@@                  @@          ##@$@           @$f       ");
	important.push("            >@@@~                 o@@@I                     Y@@           t#@@@   nY(  _x@#        ");
	important.push("             @@%@                  1[ .                      cY            @@@$   .(@@@$h          ");
	important.push("             W@@@z                                                         @@@@      @@I           ");
	important.push("               @@@@                                    /@@r                ?@@@       _@@          ");
	important.push("               .@@@@o                                  @@@c                X#$@        @@          ");
	important.push("                 d#@@@)                                $@                  }$#@        $@          ");
	important.push("                   @$@@@I                                                  @@#@@@f{l|[#@#.         ");
	important.push("                   @@@@#@@#I                           >@@@               p@@+Ijm@@@@@}            ");
	important.push("      @@@)        C@#[ ;U$@@%01                         . .             ]@@@>. .                   ");
	important.push("     /@ }@#      @@@m    @@@@@@@@C<                                    @@@#                        ");
	important.push("     /@( i@#m   n@#I   x@@} }@@@@@@@#@t|,                 ~@d        @@@$~                         ");
	important.push("      0@   r@@[@@Z   }@@&     .  Xq@@@@@@@@@@n             YB     |&$$@                            ");
	important.push("      (@U    #@@@   ?@@              iz#@@@$@@|  @@@@@&&nr?    )@@@@$ .                            ");
	important.push("       $@     J;   @@@+                  C@$@;   $@@@@@@@@@@$@@#@@@                                ");
	important.push("       @#        {$@@            ?t%dzi  @@#u   Y@@   -X0@@@$@@Y                                   ");
	important.push("        o@-    Y@@@ .            @$#$@@@@@$    ;$@?                                                ");
	important.push("         @@  +#@@[               @@@v##@@@Z    @@r                                                 ");
	important.push("         i@@@@@@[                #@$#  n}      $@                                                  ");
	important.push("           @@#@                  z@#@         @@x                                                  ");
	important.push("                                  }@@@       @@]                                                   ");
	important.push("                                   @@#B     >@@                                                    ");
	important.push("                                   +@@@    }@@z                                                    ");
	important.push("                                     @$@-  @@(                                                     ");
	important.push("                                      L@$@@%;                                                      ");
	important.push("                                      .U@@@.                                                   ");
	
	function renderDB(data,numLines){
		if(numLines < data.length){
			setTimeout(function() {
				$("#DBPre").text(assembleDB(data,numLines));
				renderDB(data,numLines + 1);
			}, 100);
		} else {
			DBdone = true;
		}
	}
	
	function assembleDB(data, numLines){
		var DB = "";
		for(var i = 0; i < numLines; i++){
			DB += data[i] + "\n";
		}
		return DB;
	}
