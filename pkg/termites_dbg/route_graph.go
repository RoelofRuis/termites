package termites_dbg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/RoelofRuis/termites/pkg/termites"

	"github.com/goccy/go-graphviz"
)

// TODO: split into smaller files and clean up

type graphWriter struct {
	rootDir      string
	writeDotFile bool
	version      int
}

type connectionEdge struct {
	FromActor  termites.NodeId
	FromPort   termites.OutPortId
	ToActor    *termites.NodeId
	ToPort     *termites.InPortId
	ViaAdapter *termites.AdapterRef
}

func (w *graphWriter) saveRoutingGraph(nodes []termites.NodeRef) string {

	var connections []connectionEdge

	for _, ref := range nodes {
		for _, out := range ref.OutPorts {
			var fromId termites.NodeId
			for _, ref := range nodes {
				if _, has := ref.OutPorts[out.Id]; has {
					fromId = ref.Id
					break
				}
			}

			for _, connection := range out.Connections {
				var toId *termites.NodeId = nil
				var toPort *termites.InPortId = nil
				if connection.In != nil {
					for _, ref := range nodes {
						if _, has := ref.InPorts[connection.In.Id]; has {
							toId = &ref.Id
							break
						}
					}
					toPort = &connection.In.Id
				}
				connections = append(connections, connectionEdge{
					FromActor:  fromId,
					FromPort:   out.Id,
					ToActor:    toId,
					ToPort:     toPort,
					ViaAdapter: connection.Adapter,
				})
			}
		}
	}

	g := routingGraph{
		name:    "routing",
		mclimit: 10.0,
		ranksep: 0.6,
		nodesep: 0.6,
	}

	connectedIns := make(map[termites.InPortId]bool)
	connectedOuts := make(map[termites.OutPortId]bool)
	for _, conn := range connections {
		connectedOuts[conn.FromPort] = true
		if conn.ToPort != nil {
			connectedIns[*conn.ToPort] = true
		}
	}

	actorRef := make(map[termites.NodeId]string)
	inPortRefs := make(map[termites.InPortId]string)
	outPortRefs := make(map[termites.OutPortId]string)
	portIdx := 0
	for i, actor := range nodes {
		nodeInPorts := make(map[string]visualizerPort)
		for inId, in := range actor.InPorts {
			ref := fmt.Sprintf("inPort_%d", portIdx)
			portIdx++
			inPortRefs[inId] = ref
			color := "gray"
			if _, isConnected := connectedIns[inId]; isConnected {
				color = "black"
			}
			nodeInPorts[ref] = visualizerPort{
				name:  in.Name,
				color: color,
			}
		}

		nodeOutPorts := make(map[string]visualizerPort)
		for outId, out := range actor.OutPorts {
			ref := fmt.Sprintf("outPort_%d", portIdx)
			portIdx++
			outPortRefs[outId] = ref
			color := "gray"
			if _, isConnected := connectedOuts[outId]; isConnected {
				color = "black"
			}
			nodeOutPorts[ref] = visualizerPort{
				name:  out.Name,
				color: color,
			}
		}

		ref := fmt.Sprintf("node_%d", i)
		actorRef[actor.Id] = ref

		g.addNode(visualizerNode{
			name:      actor.Name,
			fontcolor: "black",
			ref:       ref,
			in:        nodeInPorts,
			out:       nodeOutPorts,
		})
	}

	for ci, conn := range connections {
		outNodeRef := actorRef[conn.FromActor]
		outPortRef := outPortRefs[conn.FromPort]

		if conn.ViaAdapter != nil {
			adapterRef := fmt.Sprintf("adapter_%d", ci)
			g.addNode(visualizerNode{
				ref:   adapterRef,
				name:  conn.ViaAdapter.Name,
				shape: "rect",
			})

			g.addEdge(edge{
				color:     "black",
				style:     "solid",
				fromNode:  outNodeRef,
				fromPort:  outPortRef,
				toNode:    adapterRef,
				toPort:    "",
				arrowhead: "none",
			})

			outNodeRef = adapterRef
			outPortRef = ""
		}

		if conn.ToActor == nil || conn.ToPort == nil {
			continue
		}

		inNodeRef := actorRef[*conn.ToActor]
		inPortRef := inPortRefs[*conn.ToPort]

		edgeStyle := "solid"
		edgeColor := "black"
		if inNodeRef == "" || inPortRef == "" {
			// component to which this edge connects is not yet connected itself!
			inNodeRef = fmt.Sprintf("unconnected_%d", ci)
			edgeStyle = "dashed"
			edgeColor = "red"
			g.addNode(visualizerNode{
				ref:       inNodeRef,
				name:      "???",
				color:     "red",
				fontcolor: "red",
			})
		}

		g.addEdge(edge{
			color:    edgeColor,
			style:    edgeStyle,
			fromNode: outNodeRef,
			fromPort: outPortRef,
			toNode:   inNodeRef,
			toPort:   inPortRef,
		})
	}

	err, path := w.write(g)
	if err != nil {
		log.Printf("error creating routing graph: %v", err)
		return ""
	}
	return path
}

func (w *graphWriter) write(g routingGraph) (error, string) {
	buf := bytes.Buffer{}
	buf.WriteString(g.String())

	graph, err := graphviz.ParseBytes(buf.Bytes())
	if err != nil {
		dotPath := filepath.Join(w.rootDir, "routing-err.dot")
		_ = ioutil.WriteFile(dotPath, buf.Bytes(), 755)
		return err, ""
	}

	gv := graphviz.New()
	defer gv.Close()

	w.version += 1

	svgPath := filepath.Join(w.rootDir, fmt.Sprintf("routing-%d.svg", w.version))
	err = gv.RenderFilename(graph, graphviz.SVG, svgPath)
	if err != nil {
		return err, ""
	}

	if w.writeDotFile {
		dotPath := filepath.Join(w.rootDir, fmt.Sprintf("routing-%d.dot", w.version))
		err = gv.RenderFilename(graph, graphviz.XDOT, dotPath)
		if err != nil {
			return err, ""
		}
	}

	return nil, svgPath
}

type visualizerNode struct {
	ref           string
	name          string
	color         string
	fontcolor     string
	strikethrough bool
	in            map[string]visualizerPort
	out           map[string]visualizerPort
	subtext       string
	shape         string
}

type visualizerPort struct {
	name  string
	color string
}

func (n *visualizerNode) String() string {
	contents := ""
	if len(n.in) > 0 {
		contents += fmt.Sprintf("<td border='0'>%s</td>", ports(n.in, true))
	}

	subtext := ""
	if n.subtext != "" {
		subtext = fmt.Sprintf("<br/><font point-size='8'>%s</font>", n.subtext)
	}

	name := fmt.Sprintf("<b>%s</b>", n.name)
	if n.strikethrough {
		name = fmt.Sprintf("<s>%s</s>", name)
	}

	contents += fmt.Sprintf("<td border='0' port='default'>%s%s</td>", name, subtext)
	if len(n.out) > 0 {
		contents += fmt.Sprintf("<td border='0'>%s</td>", ports(n.out, false))
	}

	color := "black"
	if n.color != "" {
		color = n.color
	}

	fontcolor := "black"
	if n.fontcolor != "" {
		fontcolor = n.fontcolor
	}

	shape := "none"
	outline := "1"
	if n.shape != "" {
		shape = n.shape
		outline = "0"
	}

	return fmt.Sprintf(`
%s [
	color=%s
	fontcolor=%s
    shape=%s
	label=<
  <table cellborder='%s' border='%s' cellspacing='0' style='rounded'>
    <tr>
      %s
    </tr>
  </table>
>];
`, n.ref, color, fontcolor, shape, outline, outline, contents)
}

func ports(m map[string]visualizerPort, in bool) string {
	side := "r"
	if !in {
		side = "l"
	}
	rows := ""
	total := len(m) - 1
	i := 0
	for portRef, p := range m {
		var sides string
		if i == 0 {
			if total == 0 {
				sides = side
			} else {
				sides = side + "b"
			}
		} else if i == total {
			sides = side + "t"
		} else {
			sides = side + "bt"
		}
		i++
		rows += fmt.Sprintf(
			"<tr><td sides='%s' port='%s'><font color='%s'>%s</font></td></tr>\n",
			sides,
			portRef,
			p.color,
			p.name,
		)
	}
	return fmt.Sprintf("<table border='0' cellborder='1' cellspacing='0'>\n%s\n</table>\n", rows)
}

type edge struct {
	color     string
	style     string
	fromNode  string
	fromPort  string
	toNode    string
	toPort    string
	arrowhead string
}

func (e *edge) string() string {
	from := e.fromNode
	if e.fromPort != "" {
		from += ":" + e.fromPort
	}
	to := e.toNode
	if e.toPort != "" {
		to += ":" + e.toPort
	}
	color := "black"
	if e.color != "" {
		color = e.color
	}

	style := "solid"
	if e.style != "" {
		style = e.style
	}

	arrowhead := "normal"
	if e.arrowhead != "" {
		arrowhead = e.arrowhead
	}

	return fmt.Sprintf("%s->%s [color=%s style=%s arrowhead=%s];\n", from, to, color, style, arrowhead)
}

type routingGraph struct {
	name    string
	mclimit float64
	nodesep float64
	ranksep float64
	nodes   []visualizerNode
	edges   []edge
}

func (g *routingGraph) String() string {
	nodestring := ""
	for _, n := range g.nodes {
		nodestring += n.String()
	}
	edgestring := ""
	for _, e := range g.edges {
		edgestring += e.string()
	}

	return fmt.Sprintf(`
digraph %s {
    graph[rankdir=LR searchsize=100 rank=min mclimit=%.1f nodesep=%.2f ranksep=%.2f];
	node [shape=plaintext fontname="Courier" fontsize="10"];
	%s
	%s
}
`, g.name, g.mclimit, g.nodesep, g.ranksep, nodestring, edgestring)
}

func (g *routingGraph) addNode(n visualizerNode) {
	g.nodes = append(g.nodes, n)
}

func (g *routingGraph) addEdge(e edge) {
	g.edges = append(g.edges, e)
}
