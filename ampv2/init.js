var actionModels = {};

// TreeParser is based off of the c3d spec, and each function assumes that non-leaf 
// nodes has action models that adheres to the spec, and that the contracts does too. 
// This means each node must have a c3d compliant way of getting their children.
function TreeParser(){
	// Recursive. Gets a node based on the path, shaving off the first path element each time.
	// The address of the node model has been set.
	function getNode(path, startingNode) {
		var node = startingNode;
		while(path.length > 1){
			var childAddr = node.getchild(path[0]);
			node = GetNode(childAddr)
			path = path.slice[1,path.length];
		}
		return node;
	};
	
	// Get the gendoug model.
	function getGenDoug() {
		return GetContract(GENDOUG);
	};
	
	// Run
	this.run = function(path,cmd,params){
		if(!(instanceof path Array) || typeof cmd !== "string"){
			return null;
		}
		var genDoug = getGenDoug();
		if(commands === null || commands.length === 0){
			return genDoug[cmd](params);
		}
		return getNode(path,genDoug)[cmd](params);
	}
};

var treeParser = new TreeParser();

GetModelFromHash = function(modelHash){
	return actionModels[modelHash];
}

GetContract = function(contractAddress){
	var modelName = GetStorageAt(contractAddress,"0x19");  
	var modelHash = GetModelHash(modelName);
	var model = GetModelFromHash(modelHash);
	model.setAddress(contractAddress);
	return model;
}