
var dataPoints = 200;

// Heap data
var hsDraw = true;
var hsDataMin = Math.pow(2,32);
var hsDataMax = -1
var hsUpdate = true;

// Heap bytes allocated.
var hsAlloc = [];

//Heap bytes used.
var hsUsed = [];


initRPCEventHandlers = function(){
	// Done when clicking the compile button at editor.
	prof.RPCEventHandlers["MemStats"] = updateMemStats;
}

updateMemStats = function(result){
	
	if (hsAlloc.length === dataPoints){
		hsAlloc = hsAlloc.slice(1);
	}
	var hAlloc = result.HeapAlloc;
	hsAlloc.push(hAlloc);
	
	$('#MDAlloc').html(fmtMB(result.Alloc));
	$('#MDTotalAlloc').html(fmtMB(result.TotalAlloc));
	$('#MDSys').html(fmtMB(result.Sys));
	$('#MDLookups').html(fmtDiv(result.Lookups));
	$('#MDMallocs').html(fmtDiv(result.Mallocs));
	$('#MDFrees').html(fmtDiv(result.Frees));
	
	$('#MDHeapAlloc').html(fmtMB(result.HeapAlloc));
	$('#MDHeapSys').html(fmtMB(result.HeapSys));
	$('#MDHeapIdle').html(fmtMB(result.HeapIdle));
	$('#MDHeapInuse').html(fmtDiv(result.HeapInuse));
	$('#MDHeapReleased').html(fmtDiv(result.HeapReleased));
	$('#MDHeapObjects').html(fmtDiv(result.HeapObjects));
	
	$('#MDNextGC').html(result.NextGC);
	$('#MDLastGC').html(fmtTime(result.LastGC));
	$('#MDPauseTotalNs').html(fmtMS(result.PauseTotalNs));
	$('#MDNumGC').html(fmtDiv(result.NumGC));
	$('#MDEnableGC').html(result.EnableGC.toString());
	$('#MDDebugGC').html(result.DebugGC.toString());
	
	// If we are drawing a graph.
	if(hsDraw){
		
		var resAlloc = [];
		for (var j = 0; j < hsAlloc.length; ++j) {
			resAlloc.push([j, hsAlloc[j]]);
		}
		
		var data = getGraphData(resAlloc);
		
		memGraph.setData(data);
		
		if(hsDataMin > hAlloc){
			hsDataMin = hAlloc;
		}
		
		if(hsDataMax < hAlloc){
			hsDataMax = hAlloc;
		}
		
		memGraph.getYAxes()[0].options.max = hsDataMax*1.1;
		memGraph.getYAxes()[0].options.min = hsDataMin*0.9;
		memGraph.setupGrid();
		memGraph.draw();
	}
	
	if(hsUpdate){
		setTimeout(function(){prof.MemStats()},5000);
	}
}

function get2_10keyValueDiv(key, value){
	return '<div class="row"><div class="col-sm-2"><b>' + key.toString() + '</b></div><div class="col-sm-10">' + value.toString() + '</div></div>';
}

function get12_2r_keyValueDiv(key, value){
	return '<div class="row"><b>' + key.toString() + '</b></div><div class="row">' + value.toString() + '</div>';
}

function fmtMB(value){
	var val = value / 1000000;
	return val.toFixed(2) + " MB";
}

function fmtTime(value){
	value = Math.floor(value / 1000000);
	var date = new Date(value);
	var ts = date.toString();
	return ts.substring(0,ts.length - 16);
}

function fmtMS(value){
	value = Math.floor(value / 1000000);
	return fmtDiv(value) + " ms";
}

function fmtDiv(value){
	var str = value.toString(10);
	var len = str.length;
	var trips = Math.floor(len / 3);
	var rest = len % 3;
	var newStr = "";
	var ctr = 0;
	if(rest != 0){
		newStr = newStr + str.substring(0,rest); 
		if(trips != 0) {
			newStr = newStr + "'";
		}
		ctr += rest;
	}
	for(var i = 0; i < trips; i++){
		newStr = newStr + str.substring(ctr,ctr + 3); 
		if(i != trips - 1) {
			newStr = newStr + "'";
		}
		ctr += 3;
	}
	return newStr;
}