{
	"dependencies" : ["Model"],
	"name" : "GenDoug",
	
	"data" : {
		// Global
		"offset"  		: "0x10000", // The spacing used to avoid
										// collision of name spaces
										// (contracts, permissions and
										// variables)
		"GFzeros" 		: "0x1000000", // Names used buts end in at
										// least 3 empty bytes in order
										// to partition name spaces
		"colavd" 		: "0x100", // This is used to partition user
									// address space
		// C3D
		"BAindicator" 	: "0x10",
		"BAdmpointer" 	: "0x11",
		"BAUIpointer" 	: "0x12",
		"BAblob" 		: "0x13",
		"BAparent"		: "0x14",
		"BAowner"  		: "0x15",
		"BAcreator"		: "0x16",
		"BAtime"		: "0x17",
		"BAbehaviour" 	: "0x18",
		"BALLstart"		: "0x19",
		// LinkedList
		"headslot" 		: this["BALLstart"], // This will keep
												// LLstart at newest
												// element
		"tailslot"  	: "0x20",
		"countslot" 	: "0x22",
		"prowslot"  	: "0x23",
		"pbitslot"  	: "0x24",
	},
	
	"actions" : {
		// Getters (TODO add closures. most of these should be hidden)
		
		// GenDOUG
		"nextslot" : function(params) {
			return Add(params,"2");
		},
		
		"prevslot" : function(params) {
			return Add(params,"1");
		},
		
		"typeslot" : function(params) {
			return Add(params,"3");
		},
		
		"behaviourslot" : function(params) {
			return Add(params,"4");
		},
		
		"dataslot" : function(params) {
			return Add(params,"5");
		},
		
		"modelslot" : function(params) {
			return Add(params,"6");
		},
		
		"UIslot" : function(params) {
			return Add(params,"7");
		},
		
		"timeslot" : function(params) {
			return Add(params,"8");
		},
		
		"permname" : function(params) {
			return Add(params, _data.offset)
		},	
		
		// LinkedList
		"nextlink" : function(params) {
			return GetStorageAt(_address,Add(params,"2"));
		},
		
		"prevlink" : function(params) {
			return GetStorageAt(_address,Add(params,"1"));
		},
		
		"getParent" : 
		
		"getname" : function(params) {
			if(typeof params !== "string"){
				console.log("GenDoug.getname: Params should be a string.");
				return null;
			}
			// Make sure the name is an actual name.
			if (!IsZero(Mod(params,_data.GFzeros) ) ){
				return null;
			}
			return GetStorageAt(_address,params);
		},
		
		"checkperm" : function(params) {
			if(!(instanceof params Array) || params.length !== 2 || typeof params[0] !== "string" 
				|| typeof params[1] !== "string"){
				return null;
			}
			
			var name = params[0];
			var target = params[1];
			var permval = GetStorageAt(_address, permname(params[1]))
			
			if (!IsZero(Mod(name, _data.GFzeros)) ){
				return null;
			}
			
			if(permval === null && name !== "doug") {
				return null;
			}
			
			var sPos = Mod(permval,"256");
			var rPer = Mod(Div(permval,"255") _data.colavd);
			var temp = GetStorageAt(_address,Add(Mul(target, _data.colavd),rPer));
			var temp2 = Div(temp,Exp("2",sPos));
			return Mod(temp2,"0x10");
		},
		
		"register" : function(params) {
			// TODO Decide how to handle errors and return them. Since we
			// can't get return values in a
			// simple way yet (?), and at this point 'register' only returns
			// 0 or 1, I add
			// some stuff here for now.
			if(!(params instanceof Array) || params.length != 2){
				console.log("GenDoug.register: Param length != 2 (should be: ['regname', '0xaddress'] (string,string)");
				return false;
			}
			if(typeof params[0] !== "string" || params[1] !== "string"){
				console.log("GenDoug.register: Param length != 2 (should be: ['regname', '0xaddress'] (string,string)");
				return false;
			}
			Transact(_address,["register"].concat(params));
			return true;
		}
	},

}