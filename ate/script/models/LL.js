LL = {
	"name" : "LinkedList"

	//Structure
	"CTS" : function(name, key){
		return Add(Monk.StdVar.Vari(name), Add(Mul(Mod(key, Exp("0x100", "20")), Exp("0x100", "3")), Exp("0x100","2"));
	},
	"CTK" : function(slot){
		return Mod(Div(slot, Exp("0x100","3")), Exp("0x100","20"));
	},

	"TailSlot" : function(name){
		return Monk.LLLL.TailSlot(Monk.StdVar.Vari(name));
	},
	"HeadSlot" : function(name){
		return Monk.LLLL.HeadSlot(Monk.StdVar.Vari(name));
	},
	"LenSlot" : function(name){
		return Monk.LLLL.LenSlot(Monk.StdVar.Vari(name));
	},

	"MainSlot" : function(name, key){
		return Monk.LLLL.MainSlot(this.CTS(name, key));
	},
	"PrevSlot" : function(name, key){
		return Monk.LLLL.Prevlot(this.CTS(name, key));
	},
	"NextSlot" : function(name, key){
		return Monk.LLLL.NextSlot(this.CTS(name, key));
	},

	//Gets
	"GetTail" : function(addr, name){
		return GetStorageAt(addr, this.TailSlot(name));
	},
	"GetHead" : function(addr, name){
		return GetStorageAt(addr, this.HeadSlot(name));
	},
	"GetLen"  : function(addr, name){
		return GetStorageAt(addr, this.LenSlot(name));
	},

	"GetMain" : function(addr, name, key){
		return GetStorageAt(addr, this.MainSlot(name, key));
	},
	"GetPrev" : function(addr, name, key){
		return GetStorageAt(addr, this.PrevSlot(name, key));
	},
	"GetNext" : function(addr, name, key){
		return GetStorageAt(addr, this.NextSlot(name, key));
	},

	//Gets the whole list. Note the separate function which gets the keys
	"GetList" : function(addr, name){
		var list = [];
		var current = this.GetTail(addr, name);
		while(!IsZero(current)){
			list.push(this.GetMain(addr, current));
			current = this.GetNext(addr, current);
		}

		return list;
	},

	"GetKeys" : function(addr, name){
		var keys = [];
		var current = this.GetTail(addr, name);
		while(!IsZero(current)){
			list.push(this.CTK(current));
			current = this.GetNext(addr, current);
		}

		return keys;
	},
}