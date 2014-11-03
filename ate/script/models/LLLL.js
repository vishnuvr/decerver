LLLL = {
	"name" : "LLLinkedList"

	//Constants
	"TailSlotOffset"  : "0"
	"HeadSlotOffset"  : "1"
	"LenSlotOffset"   : "2"

	"LLLLSlotSize" 	  : "3"

	"EntryMainOffset" : "0"
	"EntryPrevOffset" : "1"
	"EntryNextOffset" : "2"

	//Structure
	"TailSlot" : function(base){
		return Add(base, this.TailSlotOffset);
	},
	"HeadSlot" : function(base){
		return Add(base, this.HeadSlotOffset);
	},
	"LenSlot" : function(base){
		return Add(base, this.LenSlotOffset);
	},

	"MainSlot" : function(slot){
		return Add(slot, this.EntryMainOffset);
	},
	"PrevSlot" : function(slot){
		return Add(slot, this.EntryPrevOffset);
	},
	"NextSlot" : function(slot){
		return Add(slot, this.EntryNextOffset);
	},

	//Gets
	"GetTail" : function(addr, base){
		return GetStorageAt(addr, this.TailSlot(base));
	},
	"GetHead" : function(addr, base){
		return GetStorageAt(addr, this.HeadSlot(base));
	},
	"GetLen"  : function(addr, base){
		return GetStorageAt(addr, this.LenSlot(base));
	}

	"GetMain" : function(addr, slot){
		return GetStorageAt(addr, this.MainSlot(slot));
	},
	"GetPrev" : function(addr, slot){
		return GetStorageAt(addr, this.PrevSlot(slot));
	},
	"GetNext" : function(addr, slot){
		return GetStorageAt(addr, this.NextSlot(slot));
	}

}