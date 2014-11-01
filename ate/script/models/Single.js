Monk.Single = {
	"name" : "Single",

	//Structure
	"ValueSlot" : function(varname){
		return Monk.StdVarSpace.Vari(varname);
	},

	//Gets
	"value" : function(addr, varname){
		return GetStorageAt(addr, this.ValueSlot(varname));
	},
}