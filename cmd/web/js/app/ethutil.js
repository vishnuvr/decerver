	
	eth = window.eth || {};
	
	/*
     * Variable: denom
     * 
     * An object containing all the ether denominations, and
     * their value (in terms of wei).
     * 
     * Denominations:
     * Each denomination has the worth of 10^(position * 3), where
     * position is its position in the list, starting at wei = 0.
     * 
     * - wei (1)
     * - Kwei
     * - Mwei
     * - Gwei
     * - szabo (10^12)
     * - finney (10^15)
     * - *ether* (10^18)
     * - Kether
     * - Mether
     * - Gether (10^27)
     * - Tether
     * - Pether
     * - Tether (10^36)
     * - Eether
     * - Zether
     * - Yether (10^45)
     * - Nether
     * - Dether
     * - Vether (10^54)
     * - Uether (10^56)
     * 
     */
    var denom = {};

    denom.wei = "1";
    denom.Kwei = denom.wei + "000";
    denom.Mwei = denom.Kwei + "000";
    denom.Gwei = denom.Mwei + "000";
    denom.szabo = denom.Gwei + "000";
    denom.finney = denom.szabo + "000";
    denom.ether = denom.finney + "000";
    denom.Kether = denom.ether + "000";
    denom.Mether = denom.Kether + "000";
    denom.Gether = denom.Mether + "000";
    denom.Tether = denom.Gether + "000";
    denom.Pether = denom.Tether + "000";
    denom.Eether = denom.Pether + "000";
    denom.Zether = denom.Eether + "000";
    denom.Yether = denom.Zether + "000";
    denom.Nether = denom.Yether + "000";
    denom.Dether = denom.Nether + "000";
    denom.Vether = denom.Dether + "000";
    denom.Uether = denom.Vether + "000";
    
    eth.denom = denom;
    
    var denomArr = [];
    
    denomArr[0] = "wei";
    denomArr[1] = "Kwei";
    denomArr[2] = "Mwei";
    denomArr[3] = "Gwei";
    denomArr[4] = "szabo";
    denomArr[5] = "finney";
    denomArr[6] = "ether";
    denomArr[7] = "Kether";
    denomArr[8] = "Mether";
    denomArr[9] = "Gether";
    denomArr[10] = "Tether";
    denomArr[11] = "Pether";
    denomArr[12] = "Eether";
    denomArr[13] = "Zether";
    denomArr[14] = "Yether";
    denomArr[15] = "Nether";
    denomArr[16] = "Dether";
    denomArr[17] = "Vether";
    denomArr[18] = "Uether";
    

    eth.denomArr = denomArr;
    
    /*
     * Function: ConvertWei(wei, denomType, decimals)
     * 
     * Converts between wei and other denominations. This is a hacky 
     * function, because BigInteger does not handle floats.
     * 
     * The denomination names can be found in the documentation
     * of the <denom> constants.
     * 
     * Parameters: 
     * wei - The number of wei as a string (hex or dec).
     * denomtype - the type of denomination
     * decimals - the number of decimals used for the ether value.
     * 
     * Returns:
     * The amount of ether as a string.
     * 
     */
    eth.convertWei = function(wei, denomType, decimals) {
        if (eth.denom[denomType] === null) {
            throw new ReferenceError("'denomType' is not a valid denomination.");
        }
        if (typeof denomType !== "string") {
            throw new TypeError("'denomType' is not a String.")
        }
        if (!eth.isInteger(decimals)) {
            throw new TypeError("'decimals' is not an integer.")
        }
        if (decimals < 0) {
            throw new RangeError("'decimals' cannot be a negative number.");
        }
        var val;
        try {
            val = bigInt(wei);
        } catch (error) {
            throw new TypeError("'wei' is not a valid string: " + wei + ". Error: " + error)
        }
        if (eth.isNull(wei)) {
            return "0";
        }

        val = val.multiply(Math.pow(10, decimals).toString());
        val = val.divide(denom[denomType]);
        val = val.toString();
        // TODO make better.
        if(val.length > 12){
        	return val.substring(0,12) + "...";
        }
        if (val !== "0") {
            val = val.substring(0, val.length - decimals) + "." + val.substring(val.length - decimals);
        }
        return val;
    };
	
    /*
     * Function: convertWeiToBest(wei)
     * 
     * Converts between wei and the lowest denomination that
     * gives the sum a non-fractional part.
     * 
     * The denomination names can be found in the documentation
     * of the <denom> constants.
     * 
     * Parameters: 
     * wei - The number of wei as a string (hex or dec).
     * 
     * Returns:
     * An array, where [0] is the amount as a string, and
     * [1] is the denomination (as a string).
     * 
     */
    eth.convertWeiToBest = function(wei) {
        var val;
        if(wei === ""){
        	return ["0","wei"];
        }
        try {
            val = bigInt(wei);
        } catch (error) {
            throw new TypeError("'wei' is not a valid hex or decimal string: " + wei)
        }
        if (eth.isNull(wei)) {
            return ["0","wei"];
        }

        var valStr = val.toString();
        var zeroes = 0;
        for(var i = valStr.length - 1; i > 0; i--){
        	if(valStr.charAt(i) == "0"){
        		zeroes++;
        	} else {
        		break;
        	}
        }
        
        var rest = zeroes % 3;
        
        var newValStr = valStr.substring(0,valStr.length - zeroes + rest);
        var denom = window.eth.getDenomByNumZeroes(zeroes);
        
        return [newValStr,denom];
    };
    
    eth.getDenomByNumZeroes = function(numZeroes){
    	if (!eth.isInteger(numZeroes)) {
            throw new TypeError("'numZeroes' is not an integer.")
        }
        if (numZeroes < 0) {
            throw new RangeError("'numZeroes' cannot be a negative number.");
        }
        if (numZeroes > 55) {
            throw new RangeError("'numZeroes' is larger then 55.");
        }
        return eth.denomArr[Math.floor(numZeroes/3)];
    }
    
    /**
	 * Takes a Number argument and turn into a date string.
	 * 
	 * @param {Number}
	 *            ts A javascript Number.
	 * @returns {String} A date string.
	 */
	eth.timestampToDate = function(ts) {
		
		if (!eth.isInteger(ts)) {
            throw new TypeError("'ts' is not an integer.")
        }
        if (ts < 0) {
            throw new RangeError("'ts' cannot be a negative number.");
        }
		
        if(ts === 0){
        	return "-";
        }
		
	    var date = new Date(ts * 1000);
	    /*
		 * var year = date.getFullYear().toString(); var month = (date.getMonth() +
		 * 1).toString(); var day = date.getDate().toString();
		 * 
		 * var hours = date.getHours().toString(); var minutes =
		 * date.getMinutes().toString(); var seconds = date.getSeconds().toString();
		 * 
		 * if (minutes.length === 1) { minutes = "0" + minutes; }
		 * 
		 * if (seconds.length === 1) { seconds = "0" + seconds; }
		 * 
		 * var dateString = year + "/" + month + "/" + day + ", " + hours + ":" +
		 * minutes + ":" + seconds;
		 */
		var dateString = date.toString();
	    return dateString;// + " (UNIX EPOCH: " + ts +")";
	};
	
	/*
     * Function: isNumber(n)
     * Checks if the variable n is a (finite) number.
     *
     * Parameters: 
     * n - The number.
     * 
     * Returns:
     * 'true' if n is a number, otherwise 'false'.
     */
    eth.isNumber = function(n) {
        return (typeof n === "number" && !isNaN(n) && isFinite(n));
    };

    /*
     * Function: isInteger(n)
     * Checks if the variable n is an integer.
     *
     * Parameters: 
     * n - The number.
     * 
     * Returns:
     * 'true' if n is an integer, otherwise 'false'.
     */
    eth.isInteger = function(n) {
        return (eth.isNumber(n) && n % 1 === 0);
    };
	
	// Checks if a number-string is null
    eth.isNull = function(val) {
        return (val === "0x" || val === "0" || val === "0x0") ? true : false;
    };
    

    eth.hexStringToBytes = function (hexString){
    	
    	if(typeof hexString !== "string"){
    		throw new TypeError("Param hexString is not a string.");
    	}
    	
    	if(hexString.length % 2 !== 0){
    		throw new Error("Param hexString can not be of odd length.");
    	}
    	hexString = hexString.trim();
    	console.log(hexString);
    	
    	var bytes = [];
    	var bt;
    	for (var i = 0; i < hexString.length; i += 2) {
    		// Unhandled exception...
    	    bytes.push(parseInt(hexString.substr(i, 2),16));
    	}
    	
    	return bytes;
    }
    
    eth.isStringAlphaNum = function(str){
    	return /^[\x00-\x7F]*$/.test(str);
    }
    
    eth.isStringNum = function(str){
    	return eth.isStringHex(str) || eth.isStringDec(str);
    }
    
    eth.isStringDecNum = function(str){
    	return /[0-9]/.test( str ); 
    }
    
    eth.isStringHexNum = function(str){
    	return /0x?[a-f0-9]/.test( str ); 
    }
    
    eth.stringForDisplay = function(str){
    	
    	
    	return quote(str);
		//var newstr = '"0x';
		//for (var j = 0; j < str.length; j++) {
		//	newstr += str.charCodeAt(j).toString(16);
		//}
		//newstr += '"';
		//return newstr;
		
    }
	
    var escapable = /[\\\"\x00-\x1f\x7f-\uffff]/g,
    meta = {    // table of character substitutions
        '\b': '\\b',
        '\t': '\\t',
        '\n': '\\n',
        '\f': '\\f',
        '\r': '\\r',
        '"' : '\\"',
        '\\': '\\\\'
    };

function quote(string) {

//If the string contains no control characters, no quote characters, and no
//backslash characters, then we can safely slap some quotes around it.
//Otherwise we must also replace the offending characters with safe escape
//sequences.

    escapable.lastIndex = 0;
    return escapable.test(string) ?
        '"' + string.replace(escapable, function (a) {
            var c = meta[a];
            return typeof c === 'string' ? c :
                '\\u' + ('0000' + a.charCodeAt(0).toString(16)).slice(-4);
        }) + '"' :
        '"' + string + '"';
}
    
	window.eth = eth;
