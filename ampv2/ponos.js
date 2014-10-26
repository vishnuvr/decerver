// Genesis DOUG	
Model = {
	"name" : "Ponos",
	"address" : GENDOUG,
	"data" : {
		"theBeef" : "0x30",
	},
	
	"deadbeef" : function(params) {
		var beef = GetStorageAt(this.address,this.data["theBeef"]);
		if(beef !== "deadbeef") {
			return "Someone took the beef!"
		}
		return beef;
	},
	
};