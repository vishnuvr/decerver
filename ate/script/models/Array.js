Monk.Array = {
	"name" : "Array"
 
	//Structure
	"CTS" : function(name, key){
		return Add(Monk.StdVar.Vari(name), Add(Mul(Mod(key, Exp("0x100", "20")), Exp("0x100", "3")), Exp("0x100","2"));
	},
	"CTK" : function(slot){
		return Mod(Div(slot, Exp("0x100","3")), Exp("0x100","20"));
	},

	"ESizeSlot" : function(name){
		return Monk.LLArray.ESizeSlot(Monk.StdVar.Vari(name));
	},
	"MaxESlot" : function(key){
		return Monk.LLArray.MaxESlot(this.CTS(name, key));
	},
	"StartSlot" : function(key){
		return Monk.LLArray.StartSlot(this.CTS(name, key));
	},

	//Gets
	"GetESize" : function(addr, name){
		return Monk.LLArray.GetESize(addr, Monk.StdVar.VarBase(Monk.StdVar.Vari(name)));
	},
	
	"GetMaxE" : function(addr, name, key){
		return Monk.LLArray.GetMaxE(addr, this.CTS(name, key));
	},

	"GetElement" : function(addr, name, key, index){
		return Monk.LLArray.GetElement(addr, Monk.StdVar.Vari(name), this.CTS(name, key), index)
	},
}