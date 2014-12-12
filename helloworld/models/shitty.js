{
	"Contract" : "",
	"Name" : "Shitty",
		
	"CreateFile" : function(params){
		var hash = Ipfs.PushBlock(block params)
		var msg = [];
		msg.push(hash);
		Monk.Msg(this.Contract,msg);
	}
}