LLKeyValue = {
	"name" 	: "LLKeyValue",

	//Constants
	//NONE

	//Functions
	"value" : function(addr, slot, offset){
		return GetStorageAt(addr, Add(slot, offset));
	},
}