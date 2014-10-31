{
	"name" : "Node",
	"dependencies" 	: ["Model"],
	
	// Data
	"TopIndicatorSlot" 	: Add(C3DTOP,"0x1"),
	"TopModelSlot" 		: Add(C3DTOP,"0x2"),
	"TopContentSlot"	: Add(C3DTOP,"0x4"),
	"TopParentSlot"		: Add(C3DTOP,"0x5"),
	"TopOwnerSlot"		: Add(C3DTOP,"0x6"),
	"TopCreatorSlot"	: Add(C3DTOP,"0x7"),
	"TopTimeSlot"		: Add(C3DTOP,"0x8"),
	"TopBehaviorSlot"	: Add(C3DTOP,"0x9"),
	"TopTailSlot"		: Add(C3DTOP,"0xa"),
	"TopHeadSlot"		: Add(C3DTOP,"0xb"),
	"TopListLenSlot"	: Add(C3DTOP,"0xc"),
	
	// Top
	
	"indicator" : function(params) {
		return GetStorageAt(this.address,this.TopIndicatorSlot);
	},
	"model" : function(params) {
		return GetStorageAt(this.address,this.TopModelSlot);
	},
	"content" : function(params) {
		return GetStorageAt(this.address,this.TopContentSlot);
	},
	"parent" : function(params) {
		return GetStorageAt(this.address,this.TopParentSlot);
	},
	"owner" : function(params) {
		return GetStorageAt(this.address,this.TopOwnerSlot);
	},
	"creator" : function(params) {
		return GetStorageAt(this.address,this.TopCreatorSlot);
	},
	"time" : function(params) {
		return GetStorageAt(this.address,this.TopTimeSlot);
	},
	"behavior" : function(params) {
		return GetStorageAt(this.address,this.TopBehaviorSlot);
	},
	"firstChild" : function(params) {
		return GetStorageAt(this.address,this.TopTailSlot);
	},
	"lastChild" : function(params) {
		return GetStorageAt(this.address,this.TopHeadSlot);
	},
	"numChildren" : function(params) {
		return GetStorageAt(this.address,this.TopListLenSlot);
	},
	
	// Get slots
	
	"childPrevSlot" : function(params) {
		return Add(params,"1");
	},
	"childNextSlot" : function(params) {
		return Add(params,"2");
	},
	"childTypeSlot" : function(params) {
		return Add(params,"3");
	},
	"childBehaviorSlot" : function(params) {
		return Add(params,"4");
	},
	"childContentSlot" : function(params) {
		return Add(params,"5");
	},
	"childModelSlot" : function(params) {
		return Add(params,"6");
	},
	"childTimeSlot" : function(params) {
		return Add(params,"8");
	},
	
	// Get values
	"childAddress" : function(params) {
		return GetStorageAt(this.address,params);
	},
	"childPrev" : function(params) {
		return GetStorageAt(this.address,this.childPrevSlot(params));
	},
	"childNext" : function(params) {
		return GetStorageAt(this.address,this.childNextSlot(params));
	},
	"childType" : function(params) {
		return GetStorageAt(this.address,this.childTypeSlot(params));
	},
	"childBehavior" : function(params) {
		return GetStorageAt(this.address,this.childBehaviorSlot(params));
	},
	"childContent" : function(params) {
		return GetStorageAt(this.address,this.childContentSlot(params));
	},
	"childModel" : function(params) {
		return GetStorageAt(this.address,this.childModelSlot(params));
	},
	"childTime" : function(params) {
		return GetStorageAt(this.address,this.childTimeSlot(params));
	},
	
	// Get child address by id. Calls bound treeparser method GetChildAddress,
	// which iterates over all children until it finds a match. Models where O(1) 
	// access is possible should override this.
	"childAddressById" : function(params) {
		if(typeof params !== "string"){
			return "0x";
		}
		return GetChildAddress(params,this.address);
	},
	
	"childAddresses" : function(params) {
		var children = [];
		if(IsZero(this.firstChild()){
			return children;
		}
		var current = this.TopTailSlot;
		
		while (!IsZero(current)){
			children.push(this.childAddress(current));
			current = this.childNext(current);
		}
		
		return children;
	},
	
	// Get child object
	"childById" : function(params) {
		if(typeof params !== "string"){
			return "0x";
		}
		var addr = this.childAddressById(params);
		if (typeof addr === "undefined"){
			return null;
		}
		var modelName = this.childModel(params);
		return NewModel(modelName,addr);
	},
	
	"children" : function(params) {
		var children = [];
		if(IsZero(this.tail()){
			return children;
		}
		var current = this.data.TopTailSlot;
		
		while (!IsZero(current)){
			var addr = this.childAddress(current);
			if (typeof addr === "undefined"){
				return null;
			}
			var modelName = this.childModel(current);
			children.push(NewModel(modelName,addr));
			current = this.nextChild(current);
		}
		return children;
	}
	
}