package ate

import (
	"fmt"
	"github.com/eris-ltd/thelonious/monk"
	"math/big"
)

var topStr string = "0x300000000000000000000000000000000000000000000000000000000000000"

var top *big.Int = atob(topStr)

// Using ints rather then pointers for children since the parsed tree is only
// used for display purposes and it is easier to flatten/send this way.
type Node struct {
	id        string  // the (child) id of a node (how is it referenced inside parent)
	address   string  // the contract address
	parent    string  // the parent contract address
	time      string  // time (block number)
	indicator string  // indicator (88554646XY)
	model     string  // model name/hash
	content   string  // the content
	owner     string  // the owner
	creator   string  // the creator
	behavior  string  // behavior
	children  []*Node // List of children
	leaf      bool    // is this a leaf node
}

func newNode() *Node {
	node := &Node{}
	node.children = make([]*Node, 0)
	return node
}

func (n *Node) addChild(node *Node) {
	n.children = append(n.children, node)
}

type TreeParser struct {
	ethChain *monk.MonkModule
	genDoug  string
}

func NewTreeParser(ethChain *monk.MonkModule) *TreeParser {
	tp := &TreeParser{}
	tp.ethChain = ethChain
	tp.genDoug = ""
	return tp
}

// Parse a tree starting at the node with address 'nodeAddress'. The id is the
// address of a node in its parent, so it is provided to the top node here.
// Id's will not be important later, and will then be removed.
func (tp *TreeParser) ParseTree(id, nodeAddress string) (*Node, error) {
	root, err := tp.getNode(id, nodeAddress)
	if err != nil {
		return nil, err
	}
	return root, nil
}

// Add this when we have a doug of all dougs.
/*
func (tp *TreeParser) ParseTreeByPath(path []string) (*Node, error) {
	return ParseTreeByPath(path,dougOfAllDougsAddress)
}
*/

// Parse a tree. The parsing starts at the node with address 'nodeAddress', but the
// tree root becomes the contract at the end of the path. Path is a path through the tree.
//
// Assume the node at 'nodeAddress' is called X and has children with ids A and B.
// B has children C and D, and D has child E. if path is ['B','D'], then it would
// find B in X, then find D in B. The root of the returned tree would be D, so it'd
// consist of D along with its only child, E.
//
// This makes it easy to build a repo, for example, knowing the ponos, org and repo ids.
// The command would be ParseTree(["ponosId","orgId","repoId"],addressToDougOfAllDougs)
// The tree it returns would have the repo as root, and it would only have issues (and
// then the issue comments as leaves) in it.
func (tp *TreeParser) ParseTreeByPath(path []string, nodeAddress string) (*Node, error) {

	if len(path) == 0 {
		fmt.Println("Path provided to tree-parser is empty, reading from root")
		return tp.ParseTree("ponos", nodeAddress)
	}

	var id string
	addr := nodeAddress
	var err error

	for len(path) > 0 {
		id = path[0]
		addr, err = tp.GetChildAddress(id, addr)
		if err != nil {
			return nil, err
		}
		if len(path) > 1 {
			path = path[1:]
		} else {
			break
		}
	}
	return tp.ParseTree(id, addr)
}

func (tp *TreeParser) GetChildAddress(childId, nodeAddress string) (string, error) {
	// This is a test of existence too.
	idc := big.NewInt(1)
	indicator := tp.ethChain.GetStorageAt(nodeAddress, (idc.Add(idc, top)).String())
	if indicator == "0x" {
		return "", fmt.Errorf("Error: Node has no indicator. Address: %s\n", nodeAddress)
	}
	chr := indicator[len(indicator)-2]
	if chr == 'a' || chr == 'A' {
		return "", fmt.Errorf("Error: Node is AX (no children). Address: %s\n", nodeAddress)
	}
	// We got a proper c3d BX node.
	cr := big.NewInt(10)
	cr.Add(cr, top)
	// Current child
	current := tp.ethChain.GetStorageAt(nodeAddress, cr.String())

	if current != "0x" {
		current = "0x" + current
	} else {
		return "", fmt.Errorf("Error: Node has no children. Address: %s\n", nodeAddress)
	}

	for current != "0x" {
		id := current
		addr := "0x" + tp.ethChain.GetStorageAt(nodeAddress, current)
		if id == childId {
			return addr, nil
		}

		bi := big.NewInt(2)
		bi.Add(bi, atob(current))
		current = tp.ethChain.GetStorageAt(nodeAddress, bi.String())
		if current != "0x" {
			current = "0x" + current
		}
	}
	return "", fmt.Errorf("Error: Node (%s) does not have a child with Id: %s\n", nodeAddress, childId)
}

/*
func (tp *TreeParser) getGendougNode() *Node {
	gd := newNode()
	gd.id = "genDoug"
	gd.model = "GenDoug"
	gd.address = tp.genDoug
	gd.parent = ""
	gd.leaf = false
	gd.indicator = "0x88554646BA"
	gd.time = tp.ethChain.GetStorageAt(tp.genDoug,"0x17")
	current := tp.ethChain.GetStorageAt(tp.genDoug,"0x19")
	for current != "0x" {
		id := current
		addr := tp.ethChain.GetStorageAt(tp.genDoug,current)
		ch := tp.getNode(id,addr)
		ch.parent = tp.genDoug
		gd.addChild(ch)
		bi := atob(current)
		bi = bi.Add(bi,big.NewInt(2) )
		current = tp.ethChain.GetStorageAt(tp.genDoug,bi.String())
	}
	return gd
}
*/

// Get a node recursively.
func (tp *TreeParser) getNode(id, address string) (*Node, error) {

	// This is a test of existence too.
	idc := big.NewInt(1)
	indicator := tp.ethChain.GetStorageAt(address, (idc.Add(idc, top)).String())
	if indicator == "0x" {
		return nil, fmt.Errorf("Error: Node has no indicator. Address: %s\n", address)
	}

	nd := newNode()
	nd.id = id
	nd.address = address
	nd.indicator = indicator
	// Check second to last if it's an Ax or Bx type node.
	chr := nd.indicator[len(nd.indicator)-2]
	if chr == 'a' || chr == 'A' {
		nd.leaf = true
	} else {
		nd.leaf = false
	}
	// Load it up with c3d stuff
	// Model
	mod := big.NewInt(2)
	nd.model = tp.ethChain.GetStorageAt(address, (mod.Add(mod, top)).String())
	// Time
	time := big.NewInt(8)
	nd.time = tp.ethChain.GetStorageAt(address, (time.Add(time, top)).String())
	// Parent
	par := big.NewInt(5)
	nd.parent = tp.ethChain.GetStorageAt(address, (par.Add(par, top)).String())
	// Content
	cnt := big.NewInt(4)
	nd.content = tp.ethChain.GetStorageAt(address, (cnt.Add(cnt, top)).String())
	// Owner
	onr := big.NewInt(6)
	nd.owner = tp.ethChain.GetStorageAt(address, (onr.Add(onr, top)).String())
	// Creator
	cre := big.NewInt(7)
	nd.creator = tp.ethChain.GetStorageAt(address, (cre.Add(cre, top)).String())
	// Behavior
	bhv := big.NewInt(9)
	nd.behavior = tp.ethChain.GetStorageAt(address, (bhv.Add(bhv, top)).String())

	if !nd.leaf {
		cr := big.NewInt(10)
		cr.Add(cr, top)
		// Current child
		current := tp.ethChain.GetStorageAt(address, cr.String())
		fmt.Printf("Current child: %s\n", current)
		if current != "0x" {
			current = "0x" + current
		}
		for current != "0x" {
			id := current
			addr := "0x" + tp.ethChain.GetStorageAt(address, current)
			fmt.Println(addr)
			ch, err := tp.getNode(id, addr)
			if err != nil {
				return nil, err
			}
			ch.parent = address
			nd.addChild(ch)
			bi := big.NewInt(2)
			bi.Add(bi, atob(current))
			fmt.Printf("Current child: %s\n", current)
			current = tp.ethChain.GetStorageAt(address, bi.String())
			fmt.Printf("Current child 2: %s\n", current)
			if current != "0x" {
				current = "0x" + current
			}
		}
	}
	return nd, nil
}

type FlatTree struct {
	cIdx int
	Tree []*NodeProxy
}

type NodeProxy struct {
	Index     int    // Internal use
	Id        string // the (child) id of a node (how is it referenced inside parent)
	Address   string // Address of the contract
	Parent    string // the parent contract address
	Time      string // time (block number)
	Indicator string // indicator (88554646XY)
	Model     string // model name/hash
	Content   string // the content
	Owner     string // the owner
	Creator   string // the creator
	Behavior  string // behavior
	Children  []int  // List of children (as ints)

}

// Turns a tree into a list
func (tp *TreeParser) FlattenTree(tree *Node) *FlatTree {
	flatTree := &FlatTree{cIdx: 0, Tree: make([]*NodeProxy, 0)}
	tp.flattenNode(tree, flatTree)
	tp.PrintFT(flatTree)
	return flatTree
}

func (tp *TreeParser) flattenNode(node *Node, ft *FlatTree) *NodeProxy {
	np := &NodeProxy{}
	np.Children = make([]int, 0)
	np.Id = node.id
	np.Index = ft.cIdx
	np.Address = node.address
	np.Parent = node.parent
	np.Time = node.time
	np.Indicator = node.indicator
	np.Model = node.model
	np.Content = node.content
	np.Owner = node.owner
	np.Creator = node.creator
	np.Behavior = node.behavior
	ft.Tree = append(ft.Tree, np)
	ft.cIdx++
	if len(node.children) > 0 {
		for _, nd := range node.children {
			npr := tp.flattenNode(nd, ft)
			np.Children = append(np.Children, npr.Index)
		}
	}
	return np
}

func (tp *TreeParser) PrintFT(ft *FlatTree) {
	fmt.Println("Flattened tree")
	for _, np := range ft.Tree {
		fmt.Printf("Node %d:\nid = %s\naddress=%s\nparent=%s\ntime=%s\nindicator=%s\nmodel=%s\nchildren = %v\n", np.Index, np.Id, np.Address, np.Parent, np.Time, np.Indicator, np.Model, np.Children)
	}
}
