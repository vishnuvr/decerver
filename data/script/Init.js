// Stores actionmodels by their file hash.
var actionModels = {};

// TreeParser is based off of the c3d spec, and each function assumes that non-leaf 
// nodes has action models that adheres to the spec, and that the contracts does too.
function TreeParser(){
	// Recursive. Gets a node based on the path, shaving off the first path element each time.
	// The address of the node model has been set.
	function getNode(path, startingNode) {
		var node = startingNode;
		while(path.length > 1){
			var childAddr = node.childById(path[0]);
			node = GetNode(childAddr);
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

// Get a model from its hash. Address is not set.
GetModelFromHash = function(modelHash){
	return actionModels[modelHash];
}

// Gets a new model with the given name. Address is not set. 
GetModel = function(name){
	var hash = GetModelHash(name);
	return GetModelFromHash(hash);
}

// Gets a new model with the given name. Address is pre-set (ready to be used).
NewModel = function(name, address){
	var model = GetModel(name);
	model.address = address;
	return model;
}

// Get the entire subtree from a node. Passing gendoug would return the
// entire program tree.
GetSubTree = function(node){
	
}

// Breadth first search, prints node meta data.
PrintSubTree = function(tree){
	
}

var Class = {
	
	"load" : function(name,obj){
		
		// Make sure it has a name and a dependencies field. An empty field means
		// it extends model.
		if (typeof obj.name === "undefined"){
			console.log("Class has no name field");
			return null;
		}
		// Must have a proper dependencies field.
		if (typeof obj.dependencies === "undefined" || !(obj.dependencies instanceof Array)){
			console.log("Class dependency field not set (must be an array).");
			return null;
		}
		// No multiple inheritance yet.
		if (obj.dependencies.length > 1){
			// Add this when the basics are in place. Should be order independent.
			console.log("Multiple inheritance is not yet supported.");
			return null;
		}
		// Make sure all dependencies are strings.
		for(var i = 0; i < obj.dependencies.length; i++){
			if(typeof obj.dependencies[i] !== "string"){
				console.log("Dependencies must be an array of strings.");
				return null;
			}
		}
		// If no other dependencies, then the object inherits model.
		if(obj.dependencies.length === 0){
			obj.dependencies.push("Model");
		}
		
		// Used to find circular dependencies.
		var typeMap = {};
		
		// Resolve dependencies.
		this.resolveDeps(obj,typeMap);
		
		// Now add all the dependencies.
		obj.dependencies = [];
		for (var prop in typeMap) {
			obj.dependencies.push(typeMap[prop]);
		}
		
		// Stick a few class functions in there.
		this.finalize(obj);
		
		actionModels[name] = obj;
	},
	
	// Resolve dependencies recursively. Only single inheritance allowed.
	"resolveDeps" : function(obj,typeMap) {
		// This is the endpoint (Model)
		if(obj.dependencies.length == 0){
			return;
		}
		var dep = obj.dependencies[0];
		// Circular dependency.
		if(typeof typeMap[dep] !== "undefined"){
			console.log("Circular dependency found: '" + dep + "'.");
		}
		var obj2 = GetModel(dep);
		if (obj2 === null){
			console.log("Dependency not found: '" + dep + "'.");
		}
		obj2 = this.resolveDeps(obj2,typeMap);
		this.extend(obj,obj2);
	},
	
	// B extends A by adding the fields and methods of A to B. Properties of B
	// takes precedence.
	"extend" : function(objB, objA) {
		for (var prop in objA) {
			if (objA.hasOwnProperty(prop)) {
				if (!objB.hasOwnProperty(prop)) {
					objB[prop] = objA[prop];
				}
			}
		}
	}
	
	// Add some utility functions.
	"finalize" : function(obj) {
		if(typeof obj.instanceOf !== "undefined"){
			console.log("Restricted field/method name found in '" + obj.name + "': 'instanceOf'");
		}
		obj.instanceOf = function(className){
			if(typeof className !== "string"){
				console.log("instanceOf requires a string argument");
				return false;
			}
			if(obj.name === className){
				return true;
			}
			for (var i = 0; i < obj.dependencies.length; i++){
				if(obj.dependencies[i] === className){
					return true;
				}
			}
			return false;
		}
	}
	
};