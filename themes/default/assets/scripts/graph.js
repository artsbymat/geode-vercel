const container = document.getElementById("graph-container");
const graphData = JSON.parse(container.dataset.graph);
const currentPage = container.dataset.currentPage;

let currentTheme = localStorage.getItem("theme") || "light";

const drawNode = (node, ctx, globalScale) => {
  const label = node.title || node.id;
  const fontSize = 12 / globalScale;
  ctx.font = `${fontSize}px Inter, sans-serif`;

  ctx.textAlign = "center";
  ctx.textBaseline = "middle";
  ctx.fillStyle = currentTheme === "dark" ? "#fff" : "#111";
  ctx.fillText(label, node.x, node.y);
};

const getLinkColor = () =>
  currentTheme === "dark" ? "rgba(255, 255, 255, 0.4)" : "rgba(0, 0, 0, 0.4)";

const graph = ForceGraph()(container)
  .graphData(graphData)
  .nodeId("id")
  .nodeLabel("title")
  .nodeAutoColorBy("id")
  .nodeRelSize(6)
  .nodeCanvasObjectMode(() => "replace")
  .nodeCanvasObject(drawNode)
  .linkColor(getLinkColor)
  .onNodeClick((node) => {
    window.location.href = node.url;
  })
  .width(280)
  .height(250);

document.addEventListener("theme-change", (event) => {
  currentTheme = event.detail;
  graph.nodeCanvasObject(drawNode).linkColor(getLinkColor);
});
