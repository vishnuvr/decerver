	var conn = null;

	var hasStorage = false;
	
	var frSupport;
	
	var memGraph = $.plot('#MemoryGraph',getGraphData([[0,0]],[[0,0]]),{
		stack: true,
		lines: {
			show: true,
			fill: true
			
		},
		yaxis: {
			min: 0,
			max: 1,
			tickFormatter: mbFormatter,
		    tickDecimals: 2
		},
		xaxis: {
			show: false
		}
	});

	function mbFormatter(v, axis) {
		var val = v / 1000000;
		return val.toFixed(axis.tickDecimals) + " MB";
	}
	
	function getGraphData(dt) {
		var obj = [{label: "Heap Allocated", color: "rgb(110,204,255)", data: dt}];
		return obj;
	}

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
		
		initRPCEventHandlers();
		
		if (window["WebSocket"]) {
	        conn = new WebSocket("ws://localhost:3000/srpc");
	        
	        conn.onopen = function () {
	        	console.log("WebSockets connection opened.");
	        	prof.MemStats();
	        }
	        
	        conn.onclose = function(evt) {
	            console.log("WebSockets connection closed.");
	        }
	        
	        conn.onmessage = function(evt) {	
	        	prof.handleSRPC(evt.data);
	        }
	    } else {
	    	console.log("Your browser does not support WebSockets.");
	    }
		
	});
	
	$(function() {
		// Triggers		
		
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