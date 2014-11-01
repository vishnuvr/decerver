Monk.Double = {
	"name" : "Double",

	//Structure
	"ValueSlot" : function(varname){
		return Add(StdVarSpace.Vari(varname),StdVarSpace.VarSlotSize);
	},
	"ValueSlot2" : function(varname){
		return Add(this.ValueSlot(varname),1);
	},

	//Gets
	"value" : function(addr, varname){
		var values = [];
		values.push(GetStorageAt(addr, this.ValueSlot(varname)));
		values.push(GetStorageAt(addr, this.ValueSlot2(varname)));
		return values
	},

}