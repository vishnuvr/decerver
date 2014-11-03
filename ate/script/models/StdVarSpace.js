Monk.StdVarSpace = {
	"name" : "StdVarSpace",

	//Constants
	"VarSlotSize" 	: "0x5"
	"StdVarOffset" 	: "0x1"

	//Functions?
	"Vari" 	: function(varname){
		return Add(Mul(NSBase, this.StdVarOffset)),Mul(Div(SHA3(varname),Exp("0x100", "24")),Exp("0x100", "23"));
	},
	"VarBase" 	: function(varname){
		return Add(varname, this.StdVarOffset)
	},

	//Data Slots
	"VarTypeSlot"	: function(varname){
		return this.Vari(varname);
	},
	"VarNameSlot"	: function(varname){
		return Add(this.Vari(varname), 1);
	},
	"VarAddPermSlot"	: function(varname){
		return Add(this.Vari(varname), 2);
	},
	"VarRmPermSlot" 	: function(varname){
		return Add(this.Vari(varname), 3);
	},
	"VarModPermSlot"	: function(varname){
		return Add(this.Vari(varname), 4);
	},

	//Getting Variable stuff
	"type" 	: function(addr, varname){
		return GetStorageAt(addr,this.VarTypeSlot);
	},
	"name" 	: function(addr, varname){
		return GetStorageAt(addr,this.VarNameSlot);
	},
	"addperm" 	: function(addr, varname){
		return GetStorageAt(addr,this.VarAddPermSlot);
	},
	"rmperm" 	: function(addr, varname){
		return GetStorageAt(addr,this.VarRmPermSlot);
	},
	"modperm" 	: function(addr, varname){
		return GetStorageAt(addr,this.VarModPermSlot);
	},
}