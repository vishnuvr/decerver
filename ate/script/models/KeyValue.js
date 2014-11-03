Monk.KeyValue = {
	"name" 	: "KeyValue",

	//Constants
	//None

	"CTS" : function(name, key){
		return Add(Monk.StdVar.Vari(name), Add(Mul(Mod(key, Exp("0x100", "20")), Exp("0x100", "3")), Exp("0x100","2"));
	},
	"CTK" : function(slot){
		return Mod(Div(slot, Exp("0x100","3")), Exp("0x100","20"));
	},
	
	"value" : function(addr, varname, key){
		return Monk.LLKeyValue.value(addr, KVCTS(varname, key), "0")
	},

}